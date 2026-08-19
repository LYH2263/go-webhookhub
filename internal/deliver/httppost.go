package deliver

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/LYH2263/go-webhookhub/internal/clock"
	ierr "github.com/LYH2263/go-webhookhub/internal/errors"
	"github.com/LYH2263/go-webhookhub/internal/httpx"
	"github.com/LYH2263/go-webhookhub/internal/payload"
	"github.com/LYH2263/go-webhookhub/internal/sign"
)

// Poster 执行单次 HTTP POST。
type Poster struct {
	Client    *httpx.Client
	Signer    sign.Signer
	Clock     clock.Clock
	UserAgent string
}

func (p *Poster) now() time.Time {
	if p != nil && p.Clock != nil {
		return p.Clock.Now()
	}
	return time.Now().UTC()
}

// Post 发送一次。响应 Body 用 defer Close；HTTP/状态错误用 %w 包裹哨兵。
func (p *Poster) Post(ctx context.Context, in Request) Result {
	in = cloneRequest(in)
	res := Result{
		EndpointID:  in.EndpointID,
		URL:         in.URL,
		Event:       in.Event,
		DeliveryID:  in.DeliveryID,
		RequestID:   in.RequestID,
		Attempts:    in.Attempt,
		BodyBytes:   len(in.Body),
		PayloadCopy: payload.CloneBytes(in.Body),
		FinishedAt:  p.now(),
	}
	if ctx == nil {
		ctx = context.Background()
	}
	if err := ctx.Err(); err != nil {
		res.Err = err.Error()
		return res
	}
	if p == nil || p.Client == nil {
		cause := ierr.Wrap(ierr.ErrHTTP, "nil client")
		res.Cause = cause
		res.Err = cause.Error()
		return res
	}
	signer := p.Signer
	if signer == nil {
		signer = sign.Default()
	}

	ts := p.now()
	wire, err := payload.Wrap(in.DeliveryID, in.Event, ts, in.Attempt, in.Body)
	if err != nil {
		res.Err = err.Error()
		return res
	}
	sig := signer.HeaderValue(in.Secret, wire)
	res.Signed = sig != ""

	actx, cancel := httpx.WithTimeout(ctx, in.Timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(actx, http.MethodPost, in.URL, bytes.NewReader(wire))
	if err != nil {
		res.Err = ierr.WrapErr(ierr.ErrHTTP, err).Error()
		return res
	}
	ua := p.UserAgent
	if ua == "" {
		ua = p.Client.UserAgent()
	}
	payload.ApplyHubHeaders(req.Header, in.Event, in.DeliveryID, sig, ua, in.Attempt, ts, in.Headers)
	if in.RequestID != "" {
		req.Header.Set("X-Request-Id", in.RequestID)
	}

	start := time.Now()
	resp, err := p.Client.Do(req)
	res.Duration = time.Since(start)
	res.FinishedAt = p.now()
	if err != nil {
		cause := wrapHTTPErr(err)
		if errors.Is(err, context.DeadlineExceeded) || actx.Err() == context.DeadlineExceeded {
			cause = ierr.WrapErr(ierr.ErrTimeout, err)
		}
		res.Cause = cause
		res.Err = cause.Error()
		return res
	}
	DrainAndClose(resp)

	res.StatusCode = resp.StatusCode
	if httpx.Success(resp.StatusCode) {
		res.OK = true
		return res
	}
	res.OK = false
	cause := ierr.Wrapf(ierr.ErrStatus, "status %d", resp.StatusCode)
	res.Cause = cause
	res.Err = cause.Error()
	return res
}

func wrapHTTPErr(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, context.Canceled) {
		return ierr.WrapErr(ierr.ErrCanceled, err)
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return ierr.WrapErr(ierr.ErrTimeout, err)
	}
	return ierr.WrapErr(ierr.ErrHTTP, err)
}
