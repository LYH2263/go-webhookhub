package dispatch

import (
	"context"
	"errors"
	"time"

	"github.com/LYH2263/go-webhookhub/internal/deliver"
	"github.com/LYH2263/go-webhookhub/internal/httpx"
	"github.com/LYH2263/go-webhookhub/internal/retry"
)

// Fanout 按顺序投递 jobs。每轮开始检查 ctx，重试等待走 retry.Wait。
func Fanout(ctx context.Context, jobs []Job, poster *deliver.Poster, pol retry.Policy, failFast bool) ([]deliver.Result, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	pol = pol.Normalize()
	out := make([]deliver.Result, 0, len(jobs))
	var first error
	for _, job := range jobs {
		if err := ctx.Err(); err != nil {
			if first == nil {
				first = err
			}
			out = append(out, cancelResult(job, err))
			if failFast {
				return out, err
			}
			continue
		}
		res := runJob(ctx, job, poster, pol)
		out = append(out, res)
		if !res.OK && first == nil {
			if res.Cause != nil {
				first = res.Cause
			} else if res.Err != "" {
				first = errors.New(res.Err)
			}
		}
		if failFast && !res.OK {
			return out, first
		}
	}
	if first != nil {
		// 部分失败仍返回结果切片；总错误只在全部被取消时上抛。
		if errors.Is(first, context.Canceled) || errors.Is(first, context.DeadlineExceeded) {
			return out, first
		}
	}
	return out, nil
}

func runJob(ctx context.Context, job Job, poster *deliver.Poster, pol retry.Policy) deliver.Result {
	max := job.MaxAttempts
	if max < 1 {
		max = pol.MaxAttempts
	}
	if max < 1 {
		max = 1
	}
	var last deliver.Result
	start := time.Now()
	for attempt := 1; attempt <= max; attempt++ {
		if err := ctx.Err(); err != nil {
			last.Err = err.Error()
			last.OK = false
			last.Attempts = attempt - 1
			if last.Attempts < 1 {
				last.Attempts = 1
			}
			last.Duration = time.Since(start)
			last.EndpointID = job.EndpointID
			last.URL = job.URL
			last.Event = job.Event
			last.DeliveryID = job.DeliveryID
			last.RequestID = job.RequestID
			return last
		}
		req := deliver.Request{
			URL:        job.URL,
			Secret:     job.Secret,
			Headers:    job.Headers,
			Event:      job.Event,
			DeliveryID: job.DeliveryID,
			RequestID:  job.RequestID,
			Body:       job.Body,
			Timeout:    job.Timeout,
			Attempt:    attempt,
			EndpointID: job.EndpointID,
		}
		last = poster.Post(ctx, req)
		last.Attempts = attempt
		if last.OK {
			last.Duration = time.Since(start)
			return last
		}
		retryable := httpx.Retryable(last.StatusCode) || last.StatusCode == 0
		if !retryable || !pol.ShouldRetry(attempt) {
			last.Duration = time.Since(start)
			return last
		}
		if err := retry.Wait(ctx, pol.Delay(attempt)); err != nil {
			last.Err = err.Error()
			last.Duration = time.Since(start)
			return last
		}
	}
	last.Duration = time.Since(start)
	return last
}

func cancelResult(job Job, err error) deliver.Result {
	msg := "canceled"
	if err != nil {
		msg = err.Error()
	}
	return deliver.Result{
		EndpointID: job.EndpointID,
		URL:        job.URL,
		Event:      job.Event,
		DeliveryID: job.DeliveryID,
		RequestID:  job.RequestID,
		Attempts:   0,
		OK:         false,
		Err:        msg,
		BodyBytes:  len(job.Body),
	}
}
