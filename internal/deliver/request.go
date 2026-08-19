package deliver

import (
	"time"

	"github.com/LYH2263/go-webhookhub/internal/payload"
)

// Request 单次 HTTP 投递参数。Body/Secret 调用方应已拷贝；本层再拷一次。
type Request struct {
	URL        string
	Secret     []byte
	Headers    map[string]string
	Event      string
	DeliveryID string
	RequestID  string
	Body       []byte
	Timeout    time.Duration
	Attempt    int
	EndpointID string
}

func cloneRequest(r Request) Request {
	r.Secret = payload.CloneBytes(r.Secret)
	r.Body = payload.CloneBytes(r.Body)
	r.Headers = payload.CloneStringMap(r.Headers)
	return r
}
