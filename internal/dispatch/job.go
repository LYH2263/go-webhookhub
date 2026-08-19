package dispatch

import (
	"time"

	"github.com/LYH2263/go-webhookhub/internal/payload"
)

// Job 一次扇出任务。Body 与 Secret 均为独立拷贝。
type Job struct {
	DeliveryID  string
	RequestID   string
	EndpointID  string
	URL         string
	Event       string
	Body        []byte
	Secret      []byte
	Headers     map[string]string
	Timeout     time.Duration
	MaxAttempts int
}

func cloneJob(j Job) Job {
	j.Body = payload.CloneBytes(j.Body)
	j.Secret = payload.CloneBytes(j.Secret)
	j.Headers = payload.CloneStringMap(j.Headers)
	return j
}
