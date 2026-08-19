package retry

import (
	"context"
	"errors"
	"time"
)

var ErrCanceled = errors.New("webhookhub: retry wait canceled")

// Wait 等待 d，同时监听 ctx.Done()。d<=0 时仍检查一次取消。
func Wait(ctx context.Context, d time.Duration) error {
	if ctx == nil {
		ctx = context.Background()
	}
	if err := ctx.Err(); err != nil {
		return wrapCancel(err)
	}
	if d <= 0 {
		return nil
	}
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return wrapCancel(ctx.Err())
	case <-t.C:
		return nil
	}
}

func wrapCancel(err error) error {
	if err == nil {
		return ErrCanceled
	}
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return err
	}
	return err
}

// WaitFunc 可替换等待，测试可注入立即返回。
type WaitFunc func(ctx context.Context, d time.Duration) error

func DefaultWait() WaitFunc { return Wait }
