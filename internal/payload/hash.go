package payload

import (
	"crypto/sha256"
	"encoding/hex"
)

// SHA256Hex 对拷贝后的载荷做摘要，用于日志去重键。
func SHA256Hex(b []byte) string {
	sum := sha256.Sum256(CloneBytes(b))
	return hex.EncodeToString(sum[:])
}

func Fingerprint(event string, b []byte) string {
	h := sha256.New()
	_, _ = h.Write([]byte(event))
	_, _ = h.Write([]byte{0})
	_, _ = h.Write(CloneBytes(b))
	return hex.EncodeToString(h.Sum(nil))
}
