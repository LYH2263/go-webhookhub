package sign

import (
	"crypto/hmac"
	"encoding/hex"
	"strings"
)

// Verify 校验头值是否匹配 secret+body。比较使用 hmac.Equal。
func Verify(secret, body []byte, headerValue string) bool {
	want := MAC(secret, body)
	hexPart, ok := ParseHeader(headerValue)
	if !ok {
		return false
	}
	got, err := hex.DecodeString(hexPart)
	if err != nil {
		return false
	}
	return hmac.Equal(want, got)
}

// VerifyRequest 从 HeaderName 取值校验。
func VerifyRequest(secret, body []byte, header string) bool {
	return Verify(secret, body, header)
}

// CompareHex 比较两个 hex 摘要（忽略 sha256= 前缀与大小写）。
func CompareHex(a, b string) bool {
	ha, oka := ParseHeader(a)
	hb, okb := ParseHeader(b)
	if !oka || !okb {
		return strings.EqualFold(strings.TrimSpace(a), strings.TrimSpace(b))
	}
	da, err1 := hex.DecodeString(ha)
	db, err2 := hex.DecodeString(hb)
	if err1 != nil || err2 != nil {
		return false
	}
	return hmac.Equal(da, db)
}
