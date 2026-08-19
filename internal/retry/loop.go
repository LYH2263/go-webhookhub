package retry

import (
	"context"
)

// Loop 执行 fn，失败且可重试时 Wait。fn 返回 (ok, status, err)。
// ctx 取消时：每轮开始检测到取消则不再调用 fn；重试等待被取消则立即中止。
func Loop(ctx context.Context, pol Policy, wait WaitFunc, fn func(attempt int) (ok bool, status int, err error)) (attempts int, lastStatus int, lastErr error) {
	if wait == nil {
		wait = Wait
	}
	pol = pol.Normalize()
	if ctx == nil {
		ctx = context.Background()
	}
	for attempt := 1; attempt <= pol.MaxAttempts; attempt++ {
		// 每轮开始先检查取消：已取消则不再调用 fn，直接返回取消错误。
		if err := ctx.Err(); err != nil {
			return attempts, lastStatus, err
		}
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
		// 重试等待被取消则立即中止，不再进入下一轮。
		if werr := wait(ctx, d); werr != nil {
			return attempts, lastStatus, werr
		}
	}
	return attempts, lastStatus, lastErr
}
