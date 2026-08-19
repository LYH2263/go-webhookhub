package idgen

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync/atomic"
	"time"
)

var seq uint64

// New 生成 prefix + 16 字节随机 hex。
func New(prefix string) string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		n := atomic.AddUint64(&seq, 1)
		return fmt.Sprintf("%s%x%016x", prefix, time.Now().UnixNano(), n)
	}
	return prefix + hex.EncodeToString(b[:])
}

func EndpointID() string { return New("ep_") }

func DeliveryID() string { return New("del_") }

func JobID() string { return New("job_") }

func RequestID() string { return New("req_") }

// IsPrefixed 粗检。
func IsPrefixed(id, prefix string) bool {
	return len(id) > len(prefix) && id[:len(prefix)] == prefix
}
