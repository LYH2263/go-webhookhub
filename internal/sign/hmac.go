package sign

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// Signer 投递签名器。Hub.New 保证非 nil。
type Signer interface {
	HeaderName() string
	SignHex(secret, body []byte) string
	HeaderValue(secret, body []byte) string
}

const (
	HeaderName = "X-Hub-Signature-256"
	Prefix     = "sha256="
)

// HMACSHA256 默认实现。
type HMACSHA256 struct{}

func Default() Signer { return HMACSHA256{} }

func (HMACSHA256) HeaderName() string { return HeaderName }

func (HMACSHA256) SignHex(secret, body []byte) string {
	return hex.EncodeToString(MAC(secret, body))
}

func (h HMACSHA256) HeaderValue(secret, body []byte) string {
	return Prefix + h.SignHex(secret, body)
}

// MAC 计算 HMAC-SHA256 原始摘要。secret/body 只读，不保留引用。
func MAC(secret, body []byte) []byte {
	mac := hmac.New(sha256.New, secret)
	_, _ = mac.Write(body)
	return mac.Sum(nil)
}

// HexMAC 返回不带前缀的 hex。
func HexMAC(secret, body []byte) string {
	return hex.EncodeToString(MAC(secret, body))
}

// HeaderValue 完整头值 sha256=<hex>。
func HeaderValue(secret, body []byte) string {
	return Prefix + HexMAC(secret, body)
}
