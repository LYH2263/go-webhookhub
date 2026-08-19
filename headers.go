package webhookhub

// 投递请求使用的标准头。签名头格式为 sha256=<hex>。
const (
	HeaderSignature = "X-Hub-Signature-256"
	HeaderEvent     = "X-Hub-Event"
	HeaderDelivery  = "X-Hub-Delivery"
	HeaderTimestamp = "X-Hub-Timestamp"
	HeaderAttempt   = "X-Hub-Attempt"
	HeaderUA        = "User-Agent"
	DefaultUA       = "go-webhookhub/1.0"
	ContentJSON     = "application/json; charset=utf-8"
)

// SignaturePrefix HMAC-SHA256 头值前缀。
const SignaturePrefix = "sha256="
