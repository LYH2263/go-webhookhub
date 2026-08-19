package payload

import (
	"encoding/json"
	"fmt"
)

// Pretty 把 JSON 美化；非 JSON 原样拷贝返回。
func Pretty(b []byte) []byte {
	b = CloneBytes(b)
	if !json.Valid(b) {
		return b
	}
	var buf []byte
	buf, err := json.MarshalIndent(json.RawMessage(b), "", "  ")
	if err != nil {
		return b
	}
	return buf
}

// Sizeof 载荷字节数，nil 为 0。
func Sizeof(b []byte) int { return len(b) }

// Truncate 截断用于日志预览，返回新切片。
func Truncate(b []byte, n int) []byte {
	if n <= 0 {
		return nil
	}
	if len(b) <= n {
		return CloneBytes(b)
	}
	out := make([]byte, n)
	copy(out, b[:n])
	return out
}

// Preview 生成可打印预览字符串。
func Preview(b []byte, n int) string {
	t := Truncate(b, n)
	return fmt.Sprintf("%q", t)
}
