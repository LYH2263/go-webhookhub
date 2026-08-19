package httpx

import (
	"net/http"
	"strconv"
	"time"
)

// RetryAfter 解析 Retry-After 头（秒或 HTTP 日期）。解析失败返回 fallback。
func RetryAfter(h http.Header, fallback time.Duration) time.Duration {
	if h == nil {
		return fallback
	}
	v := h.Get("Retry-After")
	if v == "" {
		return fallback
	}
	if n, err := strconv.Atoi(v); err == nil {
		if n < 0 {
			return fallback
		}
		d := time.Duration(n) * time.Second
		if d > 30*time.Second {
			d = 30 * time.Second
		}
		return d
	}
	if t, err := http.ParseTime(v); err == nil {
		d := time.Until(t)
		if d < 0 {
			return 0
		}
		if d > 30*time.Second {
			return 30 * time.Second
		}
		return d
	}
	return fallback
}
