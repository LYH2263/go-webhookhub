package webhookhub

import "time"

// Endpoint 一次订阅目标。Secret / Events / Headers 在入库与返回时都会拷贝。
type Endpoint struct {
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

// EndpointView 对外只读视图。Secret 为独立拷贝，调用方改切片不会写穿库存。
type EndpointView struct {
	ID          string            `json:"id"`
	URL         string            `json:"url"`
	Events      []string          `json:"events"`
	Enabled     bool              `json:"enabled"`
	SecretLen   int               `json:"secret_len"`
	Secret      []byte            `json:"-"`
	Headers     map[string]string `json:"headers,omitempty"`
	Timeout     time.Duration     `json:"timeout_ns"`
	MaxAttempts int               `json:"max_attempts"`
	Description string            `json:"description,omitempty"`
}

// DeliveryResult 一次 Dispatch 对单个端点的最终结果（含重试）。
type DeliveryResult struct {
	EndpointID string        `json:"endpoint_id"`
	URL        string        `json:"url"`
	Event      string        `json:"event"`
	DeliveryID string        `json:"delivery_id"`
	Attempts   int           `json:"attempts"`
	StatusCode int           `json:"status_code"`
	OK         bool          `json:"ok"`
	Err        string        `json:"err,omitempty"`
	Duration   time.Duration `json:"duration_ns"`
	Signed     bool          `json:"signed"`
	BodyBytes  int           `json:"body_bytes"`
}

// DeliveryRecord 投递日志条目，RecentDeliveries 返回的是拷贝。
type DeliveryRecord struct {
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
}

// Stats 运行计数，供管理页展示。
type Stats struct {
	Endpoints     int       `json:"endpoints"`
	Enabled       int       `json:"enabled"`
	Dispatches    uint64    `json:"dispatches"`
	Deliveries    uint64    `json:"deliveries"`
	Successes     uint64    `json:"successes"`
	Failures      uint64    `json:"failures"`
	LastEvent     string    `json:"last_event,omitempty"`
	LastAt        time.Time `json:"last_at,omitempty"`
	JournalLen    int       `json:"journal_len"`
	Closed        bool      `json:"closed"`
	DefaultMaxTry int       `json:"default_max_try"`
}
