package webhookhub

import (
	"net/http"
	"time"

	"github.com/LYH2263/go-webhookhub/internal/clock"
	"github.com/LYH2263/go-webhookhub/internal/sign"
)

// Option 在 New 时注入依赖。New 结束后 signer / client 保证非 nil。
type Option func(*Hub)

func WithClock(c clock.Clock) Option {
	return func(h *Hub) {
		if c != nil {
			h.clk = c
		}
	}
}

func WithHTTPTimeout(d time.Duration) Option {
	return func(h *Hub) {
		if d > 0 {
			h.httpTimeout = d
		}
	}
}

func WithTransport(rt http.RoundTripper) Option {
	return func(h *Hub) {
		if rt != nil {
			h.transport = rt
		}
	}
}

func WithMaxAttempts(n int) Option {
	return func(h *Hub) {
		if n > 0 {
			h.maxAttempts = n
		}
	}
}

func WithBackoff(base, cap time.Duration) Option {
	return func(h *Hub) {
		if base >= 0 {
			h.backoffBase = base
		}
		if cap >= 0 {
			h.backoffCap = cap
		}
	}
}

func WithBackoffFactor(f float64) Option {
	return func(h *Hub) {
		if f >= 1 {
			h.backoffFactor = f
		}
	}
}

func WithJournalCap(n int) Option {
	return func(h *Hub) {
		if n > 0 {
			h.journalCap = n
		}
	}
}

func WithPersistPath(path string) Option {
	return func(h *Hub) { h.persistPath = path }
}

func WithJournalPath(path string) Option {
	return func(h *Hub) { h.journalPath = path }
}

func WithMaxEndpoints(n int) Option {
	return func(h *Hub) {
		if n > 0 {
			h.maxEndpoints = n
		}
	}
}

func WithUserAgent(ua string) Option {
	return func(h *Hub) {
		if ua != "" {
			h.userAgent = ua
		}
	}
}

func WithSigner(s sign.Signer) Option {
	return func(h *Hub) {
		if s != nil {
			h.signer = s
		}
	}
}

func WithAllowHTTP(v bool) Option {
	return func(h *Hub) { h.allowHTTP = v }
}

func WithFailFast(v bool) Option {
	return func(h *Hub) { h.failFast = v }
}

func WithEmptyBodyOK(v bool) Option {
	return func(h *Hub) { h.emptyBodyOK = v }
}

// EndpointSink 订阅后的持久化钩子。失败时 Subscribe 必须回滚。
type EndpointSink interface {
	SaveEndpoints(views []EndpointView) error
}

func WithEndpointSink(s EndpointSink) Option {
	return func(h *Hub) { h.sink = s }
}
