package deliver

import (
	"time"

	"github.com/LYH2263/go-webhookhub/internal/payload"
)

// Result 单次（含重试循环后的）投递结果。
type Result struct {
	EndpointID  string
	URL         string
	Event       string
	DeliveryID  string
	RequestID   string
	Attempts    int
	StatusCode  int
	OK          bool
	Err         string
	Duration    time.Duration
	Signed      bool
	BodyBytes   int
	FinishedAt  time.Time
	PayloadCopy []byte
}

func (r Result) Clone() Result {
	r.PayloadCopy = payload.CloneBytes(r.PayloadCopy)
	return r
}
