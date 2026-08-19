package webhookhub

import (
	"context"
	"errors"

	"github.com/LYH2263/go-webhookhub/internal/deliver"
	"github.com/LYH2263/go-webhookhub/internal/dispatch"
	"github.com/LYH2263/go-webhookhub/internal/journal"
	"github.com/LYH2263/go-webhookhub/internal/payload"
	"github.com/LYH2263/go-webhookhub/internal/retry"
)

// Dispatch 使用 Background context 扇出投递。
func (h *Hub) Dispatch(event string, body []byte) ([]DeliveryResult, error) {
	return h.DispatchContext(context.Background(), event, body)
}

// DispatchContext 按事件匹配端点并投递。body 立即拷贝，调用方随后改切片不影响投递。
// ctx 取消会中止尚未开始的端点以及重试等待。Close 之后返回 ErrClosed，不触碰已清空的 client。
func (h *Hub) DispatchContext(ctx context.Context, event string, body []byte) ([]DeliveryResult, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if event == "" {
		return nil, ErrInvalidEvent
	}

	h.mu.Lock()
	if err := h.checkOpenLocked(); err != nil {
		h.mu.Unlock()
		return nil, err
	}
	if !h.emptyBodyOK && len(body) == 0 {
		h.mu.Unlock()
		return nil, ErrEmptyBody
	}
	deps := h.snapshotLocked()
	if deps.client == nil {
		h.mu.Unlock()
		return nil, ErrClosed
	}
	bodyCopy := payload.CloneBytes(body)
	targets := h.reg.Match(event)
	h.dispatches++
	h.lastEvent = event
	h.lastAt = deps.clk.Now()
	h.mu.Unlock()

	jobs := dispatch.Plan(event, bodyCopy, targets, deps.policy.MaxAttempts, h.httpTimeout)
	poster := &deliver.Poster{
		Client:    deps.client,
		Signer:    deps.signer,
		Clock:     deps.clk,
		UserAgent: deps.userAgent,
	}
	raw, err := dispatch.Fanout(ctx, jobs, poster, deps.policy, deps.failFast)
	results := make([]DeliveryResult, 0, len(raw))
	var succ, fail uint64
	for _, r := range raw {
		dr := toPublicResult(r)
		results = append(results, dr)
		if r.OK {
			succ++
		} else {
			fail++
		}
		if deps.log != nil {
			deps.log.Append(toJournalRecord(r))
		}
	}

	h.mu.Lock()
	h.deliveries += uint64(len(results))
	h.successes += succ
	h.failures += fail
	h.mu.Unlock()

	if err != nil && (errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) || errors.Is(err, retry.ErrCanceled)) {
		return results, err
	}
	return results, err
}

func toPublicResult(r deliver.Result) DeliveryResult {
	return DeliveryResult{
		EndpointID: r.EndpointID,
		URL:        r.URL,
		Event:      r.Event,
		DeliveryID: r.DeliveryID,
		Attempts:   r.Attempts,
		StatusCode: r.StatusCode,
		OK:         r.OK,
		Err:        r.Err,
		Duration:   r.Duration,
		Signed:     r.Signed,
		BodyBytes:  r.BodyBytes,
	}
}

func toJournalRecord(r deliver.Result) journal.Record {
	return journal.Record{
		ID:         r.DeliveryID,
		EndpointID: r.EndpointID,
		URL:        r.URL,
		Event:      r.Event,
		StatusCode: r.StatusCode,
		OK:         r.OK,
		Attempts:   r.Attempts,
		At:         r.FinishedAt,
		Latency:    r.Duration,
		BodySize:   r.BodyBytes,
		Error:      r.Err,
		RequestID:  r.RequestID,
		Payload:    payload.CloneBytes(r.PayloadCopy),
	}
}

// RecentDeliveries 返回最近 n 条投递记录（拷贝）。n<=0 时返回全部。
func (h *Hub) RecentDeliveries(n int) []DeliveryRecord {
	h.mu.Lock()
	log := h.log
	h.mu.Unlock()
	if log == nil {
		return nil
	}
	recs := log.Recent(n)
	out := make([]DeliveryRecord, 0, len(recs))
	for _, r := range recs {
		out = append(out, DeliveryRecord{
			ID:         r.ID,
			EndpointID: r.EndpointID,
			URL:        r.URL,
			Event:      r.Event,
			StatusCode: r.StatusCode,
			OK:         r.OK,
			Attempts:   r.Attempts,
			At:         r.At,
			Latency:    r.Latency,
			BodySize:   r.BodySize,
			Error:      r.Error,
			RequestID:  r.RequestID,
		})
	}
	return out
}
