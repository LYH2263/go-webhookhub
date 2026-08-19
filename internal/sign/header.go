package sign

import (
	"net/http"
	"strings"
)

// SetHeader 把签名写入请求。空签名则删除该头。
func SetHeader(h http.Header, value string) {
	if h == nil {
		return
	}
	if value == "" {
		h.Del(HeaderName)
		return
	}
	h.Set(HeaderName, value)
}

// GetHeader 读取签名头，兼容大小写别名。
func GetHeader(h http.Header) string {
	if h == nil {
		return ""
	}
	if v := h.Get(HeaderName); v != "" {
		return v
	}
	if v := h.Get("X-Hub-Signature"); v != "" {
		return v
	}
	return ""
}

// StripPrefix 去掉 sha256= 前缀（大小写不敏感）。
func StripPrefix(v string) string {
	v = strings.TrimSpace(v)
	if len(v) >= len(Prefix) && strings.EqualFold(v[:len(Prefix)], Prefix) {
		return v[len(Prefix):]
	}
	return v
}

// ParseHeader 解析头值为 hex 部分。
func ParseHeader(v string) (hexPart string, ok bool) {
	v = strings.TrimSpace(v)
	if v == "" {
		return "", false
	}
	hexPart = StripPrefix(v)
	if hexPart == "" {
		return "", false
	}
	for _, c := range hexPart {
		if !isHex(c) {
			return "", false
		}
	}
	if len(hexPart)%2 != 0 {
		return "", false
	}
	return hexPart, true
}

func isHex(c rune) bool {
	return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
}
