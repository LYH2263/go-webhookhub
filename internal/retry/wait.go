package retry

import (
	"context"
	"errors"
	"time"
)

var ErrCanceled = errors.New("webhookhub: retry wait canceled")

// Wait 等待 d，同时监听 ctx.Done()；ctx 取消时立即返回取消错误，不睡满 d。
// d<=0 时仍检查一次取消：ctx 已取消则返回其错误，否则返回 nil。
func Wait(ctx context.Context, d time.Duration) error {
	if d <= 0 {
		return ctx.Err()
	}
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-t.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// WaitFunc 可替换等待，测试可注入立即返回。
type WaitFunc func(ctx context.Context, d time.Duration) error

func DefaultWait() WaitFunc { return Wait }
