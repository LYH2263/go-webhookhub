package journal

import (
	"time"

	"github.com/LYH2263/go-webhookhub/internal/payload"
)

// Record 投递日志。Payload 为投递时 body 的拷贝。
type Record struct {
	ID         string        `json:"id"`
	EndpointID string        `json:"endpoint_id"`
	URL        string        `json:"url"`
	Event      string        `json:"event"`
	StatusCode int           `json:"status_code"`
	OK         bool          `json:"ok"`
	Attempts   int           `json:"attempts"`
	At         time.Time     `json:"at"`
	Latency    time.Duration `json:"latency_ns"`
	BodySize   int           `json:"body_size"`
	Error      string        `json:"error,omitempty"`
	RequestID  string        `json:"request_id,omitempty"`
	Payload    []byte        `json:"payload,omitempty"`
}

func CloneRecord(r Record) Record {
	r.Payload = payload.CloneBytes(r.Payload)
	return r
}

func CloneRecords(in []Record) []Record {
	if in == nil {
		return nil
	}
	out := make([]Record, len(in))
	for i := range in {
		out[i] = CloneRecord(in[i])
	}
	return out
}
