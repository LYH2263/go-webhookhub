package payload

import (
	"net/http"
	"strconv"
	"time"
)

const (
	HeaderSignature = "X-Hub-Signature-256"
	HeaderEvent     = "X-Hub-Event"
	HeaderDelivery  = "X-Hub-Delivery"
	HeaderTimestamp = "X-Hub-Timestamp"
	HeaderAttempt   = "X-Hub-Attempt"
	HeaderUA        = "User-Agent"
	ContentType     = "Content-Type"
	ContentJSON     = "application/json; charset=utf-8"
)

// ApplyHubHeaders 写入投递标准头。extra 会再拷贝一层后合并（不覆盖标准头）。
func ApplyHubHeaders(h http.Header, event, deliveryID, signature, ua string, attempt int, ts time.Time, extra map[string]string) {
	if h == nil {
		return
	}
	h.Set(ContentType, ContentJSON)
	if ua != "" {
		h.Set(HeaderUA, ua)
	}
	if event != "" {
		h.Set(HeaderEvent, event)
	}
	if deliveryID != "" {
		h.Set(HeaderDelivery, deliveryID)
	}
	if !ts.IsZero() {
		h.Set(HeaderTimestamp, ts.UTC().Format(time.RFC3339Nano))
	}
	if attempt > 0 {
		h.Set(HeaderAttempt, strconv.Itoa(attempt))
	}
	if signature != "" {
		h.Set(HeaderSignature, signature)
	}
	for k, v := range CloneStringMap(extra) {
		if k == "" {
			continue
		}
		if h.Get(k) == "" {
			h.Set(k, v)
		}
	}
}

// CloneHeader 深拷贝 http.Header。
func CloneHeader(src http.Header) http.Header {
	if src == nil {
		return nil
	}
	out := make(http.Header, len(src))
	for k, vs := range src {
		cp := make([]string, len(vs))
		copy(cp, vs)
		out[k] = cp
	}
	return out
}
