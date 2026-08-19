package endpoint

import (
	"time"

	"github.com/LYH2263/go-webhookhub/internal/payload"
)

// Record 库存中的端点。Secret 必须是独立拷贝。
type Record struct {
	ID          string
	URL         string
	Secret      []byte
	Events      []string
	Enabled     bool
	Headers     map[string]string
	Timeout     time.Duration
	MaxAttempts int
	Description string
}

func CloneRecord(r Record) Record {
	return Record{
		ID:          r.ID,
		URL:         r.URL,
		Secret:      r.Secret,
		Events:      payload.CloneStrings(r.Events),
		Enabled:     r.Enabled,
		Headers:     payload.CloneStringMap(r.Headers),
		Timeout:     r.Timeout,
		MaxAttempts: r.MaxAttempts,
		Description: r.Description,
	}
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

// Snapshot 可 JSON 序列化的端点快照。
type Snapshot struct {
	ID          string            `json:"id"`
	URL         string            `json:"url"`
	Secret      []byte            `json:"secret"`
	Events      []string          `json:"events"`
	Enabled     bool              `json:"enabled"`
	Headers     map[string]string `json:"headers,omitempty"`
	TimeoutNS   int64             `json:"timeout_ns"`
	MaxAttempts int               `json:"max_attempts"`
	Description string            `json:"description,omitempty"`
}

func (r Record) Snapshot() Snapshot {
	return Snapshot{
		ID:          r.ID,
		URL:         r.URL,
		Secret:      payload.CloneBytes(r.Secret),
		Events:      payload.CloneStrings(r.Events),
		Enabled:     r.Enabled,
		Headers:     payload.CloneStringMap(r.Headers),
		TimeoutNS:   int64(r.Timeout),
		MaxAttempts: r.MaxAttempts,
		Description: r.Description,
	}
}

func FromSnapshot(s Snapshot) Record {
	return Record{
		ID:          s.ID,
		URL:         s.URL,
		Secret:      payload.CloneBytes(s.Secret),
		Events:      payload.CloneStrings(s.Events),
		Enabled:     s.Enabled,
		Headers:     payload.CloneStringMap(s.Headers),
		Timeout:     time.Duration(s.TimeoutNS),
		MaxAttempts: s.MaxAttempts,
		Description: s.Description,
	}
}
