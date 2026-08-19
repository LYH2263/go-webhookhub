package retry

import (
	"net"
	"os"
	"strings"
)

// RetryableStatus HTTP 状态是否值得重试。
func RetryableStatus(code int) bool {
	switch code {
	case 0:
		return true
	case 408, 425, 429:
		return true
	}
	return code >= 500 && code <= 599
}

// RetryableError 网络类错误默认可重试。
func RetryableError(err error) bool {
	if err == nil {
		return false
	}
	if ne, ok := err.(net.Error); ok && (ne.Timeout() || ne.Temporary()) {
		return true
	}
	if os.IsTimeout(err) {
		return true
	}
	msg := strings.ToLower(err.Error())
	for _, k := range []string{"connection reset", "broken pipe", "eof", "i/o timeout", "connection refused"} {
		if strings.Contains(msg, k) {
			return true
		}
	}
	return false
}

// GiveUp 不可重试的客户端错误（除 408/429）。
func GiveUp(code int) bool {
	if code >= 400 && code < 500 {
		return !RetryableStatus(code)
	}
	return false
}
