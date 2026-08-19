package dispatch

import (
	"time"

	"github.com/LYH2263/go-webhookhub/internal/endpoint"
	"github.com/LYH2263/go-webhookhub/internal/idgen"
	"github.com/LYH2263/go-webhookhub/internal/payload"
)

// Plan 把匹配到的端点编成有序 Job 列表。body 再拷一次，避免与调用方/库存共享。
func Plan(event string, body []byte, targets []endpoint.Record, defaultAttempts int, defaultTimeout time.Duration) []Job {
	body = payload.CloneBytes(body)
	if defaultAttempts < 1 {
		defaultAttempts = 1
	}
	if defaultTimeout <= 0 {
		defaultTimeout = 10 * time.Second
	}
	jobs := make([]Job, 0, len(targets))
	for _, t := range targets {
		if !t.Enabled {
			continue
		}
		maxTry := t.MaxAttempts
		if maxTry < 1 {
			maxTry = defaultAttempts
		}
		to := t.Timeout
		if to <= 0 {
			to = defaultTimeout
		}
		j := Job{
			DeliveryID:  idgen.DeliveryID(),
			RequestID:   idgen.RequestID(),
			EndpointID:  t.ID,
			URL:         t.URL,
			Event:       event,
			Body:        payload.CloneBytes(body),
			Secret:      payload.CloneBytes(t.Secret),
			Headers:     payload.CloneStringMap(t.Headers),
			Timeout:     to,
			MaxAttempts: maxTry,
		}
		jobs = append(jobs, j)
	}
	return jobs
}

// CountEnabled 统计计划中的任务数。
func CountEnabled(targets []endpoint.Record) int {
	n := 0
	for _, t := range targets {
		if t.Enabled {
			n++
		}
	}
	return n
}
