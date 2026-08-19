package httpx

import (
	"context"
	"net/http"
	"time"
)

// WithTimeout 在父 ctx 上叠加超时。timeout<=0 时原样返回父 ctx（不 cancel）。
func WithTimeout(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	if parent == nil {
		parent = context.Background()
	}
	if timeout <= 0 {
		return parent, func() {}
	}
	return context.WithTimeout(parent, timeout)
}

func IsTimeout(err error) bool {
	if err == nil {
		return false
	}
	if err == context.DeadlineExceeded {
		return true
	}
	return false
}

func MethodPOST() string { return http.MethodPost }
