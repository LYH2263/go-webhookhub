package webhookhub

import (
	"net/http"
	"sync"
	"time"

	"github.com/LYH2263/go-webhookhub/internal/clock"
	"github.com/LYH2263/go-webhookhub/internal/endpoint"
	"github.com/LYH2263/go-webhookhub/internal/httpx"
	"github.com/LYH2263/go-webhookhub/internal/journal"
	"github.com/LYH2263/go-webhookhub/internal/retry"
	"github.com/LYH2263/go-webhookhub/internal/sign"
)

const (
	defaultHTTPTimeout = 10 * time.Second
	defaultMaxAttempts = 3
	defaultBackoffBase = 50 * time.Millisecond
	defaultBackoffCap  = 2 * time.Second
	defaultFactor      = 2.0
	defaultJournalCap  = 1024
	defaultMaxEP       = 4096
)

// Hub Webhook 投递中心。零值不可用，必须经 New 构造。
type Hub struct {
	mu sync.Mutex

	closed bool
	clk    clock.Clock
	reg    *endpoint.Registry
	log    *journal.Log
	client *httpx.Client
	signer sign.Signer
	policy retry.Policy

	httpTimeout   time.Duration
	transport     http.RoundTripper
	maxAttempts   int
	backoffBase   time.Duration
	backoffCap    time.Duration
	backoffFactor float64
	journalCap    int
	maxEndpoints  int
	userAgent     string
	persistPath   string
	journalPath   string
	allowHTTP     bool
	failFast      bool
	emptyBodyOK   bool
	sink          EndpointSink

	dispatches uint64
	deliveries uint64
	successes  uint64
	failures   uint64
	lastEvent  string
	lastAt     time.Time
}

// New 构造 Hub。选项跑完后若 signer 仍为空则安装默认 HMAC-SHA256。
func New(opts ...Option) *Hub {
	h := &Hub{
		clk:           clock.Real{},
		httpTimeout:   defaultHTTPTimeout,
		maxAttempts:   defaultMaxAttempts,
		backoffBase:   defaultBackoffBase,
		backoffCap:    defaultBackoffCap,
		backoffFactor: defaultFactor,
		journalCap:    defaultJournalCap,
		maxEndpoints:  defaultMaxEP,
		userAgent:     DefaultUA,
		allowHTTP:     true,
		signer:        sign.Default(),
	}
	for _, o := range opts {
		if o != nil {
			o(h)
		}
	}
	if h.clk == nil {
		h.clk = clock.Real{}
	}
	if h.signer == nil {
		h.signer = sign.Default()
	}
	if h.maxAttempts < 1 {
		h.maxAttempts = 1
	}
	if h.journalCap < 8 {
		h.journalCap = 8
	}
	h.reg = endpoint.NewRegistry(h.maxEndpoints, h.allowHTTP)
	h.log = journal.New(h.journalCap, h.clk, h.journalPath)
	h.client = httpx.New(h.httpTimeout, h.transport, h.userAgent)
	h.policy = retry.Policy{
		MaxAttempts: h.maxAttempts,
		Base:        h.backoffBase,
		Cap:         h.backoffCap,
		Factor:      h.backoffFactor,
	}
	if h.persistPath != "" && h.sink == nil {
		h.sink = fileSinkAdapter{inner: journal.NewFileSink(h.persistPath)}
		if loaded, err := journal.LoadEndpointFile(h.persistPath); err == nil {
			for _, snap := range loaded {
				_ = h.reg.Restore(endpoint.FromSnapshot(snap))
			}
		}
	}
	return h
}

func (h *Hub) checkOpenLocked() error {
	if h.closed {
		return ErrClosed
	}
	return nil
}

func (h *Hub) snapshotLocked() dispatchDeps {
	return dispatchDeps{
		client:    h.client,
		signer:    h.signer,
		policy:    h.policy,
		clk:       h.clk,
		log:       h.log,
		failFast:  h.failFast,
		userAgent: h.userAgent,
	}
}

type dispatchDeps struct {
	client    *httpx.Client
	signer    sign.Signer
	policy    retry.Policy
	clk       clock.Clock
	log       *journal.Log
	failFast  bool
	userAgent string
}
