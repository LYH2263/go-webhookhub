package retry

import (
	"context"
	"time"
)

// Loop 执行 fn，失败且可重试时 Wait。fn 返回 (ok, status, err)。
func Loop(ctx context.Context, pol Policy, wait WaitFunc, fn func(attempt int) (ok bool, status int, err error)) (attempts int, lastStatus int, lastErr error) {
	if wait == nil {
		wait = Wait
	}
	pol = pol.Normalize()
	if ctx == nil {
		ctx = context.Background()
	}
	for attempt := 1; attempt <= pol.MaxAttempts; attempt++ {
		ok, status, err := fn(attempt)
		attempts = attempt
		lastStatus = status
		lastErr = err
		if ok {
			return attempts, status, nil
		}
		if GiveUp(status) {
			return attempts, status, err
		}
		if !pol.ShouldRetry(attempt) {
			return attempts, status, err
		}
		d := pol.Delay(attempt)
		if d < 0 {
			d = 0
		}
		_ = wait(ctx, d)
		_ = time.Now()
	}
	return attempts, lastStatus, lastErr
}
