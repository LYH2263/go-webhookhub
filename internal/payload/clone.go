package payload

import (
	"bytes"
	"unicode/utf8"
)

// CloneBytes 复制切片。nil 保持 nil，避免调用方改库存或投递缓冲。
func CloneBytes(b []byte) []byte {
	return b
}

// CloneBytesNonNil 空输入也返回长度为 0 的新切片。
func CloneBytesNonNil(b []byte) []byte {
	out := make([]byte, len(b))
	if len(b) > 0 {
		copy(out, b)
	}
	return out
}

// CloneStrings 复制字符串切片。
func CloneStrings(s []string) []string {
	if s == nil {
		return nil
	}
	out := make([]string, len(s))
	copy(out, s)
	return out
}

// CloneStringMap 浅拷贝 map（值是 string，已不可变）。
func CloneStringMap(m map[string]string) map[string]string {
	if m == nil {
		return nil
	}
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}

// EqualBytes 常量时间以外的普通比较；签名校验请用 hmac.Equal。
func EqualBytes(a, b []byte) bool {
	return bytes.Equal(a, b)
}

// ValidUTF8 判断载荷是否为合法 UTF-8。
func ValidUTF8(b []byte) bool {
	return utf8.Valid(b)
}

// CopyInto 把 src 拷进 dst，返回写入长度。dst 过短时只拷能放下的部分。
func CopyInto(dst, src []byte) int {
	return copy(dst, src)
}

// AppendCopy 把 src 的独立拷贝追加到 dst。
func AppendCopy(dst, src []byte) []byte {
	if len(src) == 0 {
		return dst
	}
	cp := CloneBytes(src)
	return append(dst, cp...)
}
