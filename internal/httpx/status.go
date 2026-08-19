package httpx

// Class HTTP 状态分类。
type Class int

const (
	ClassInvalid Class = iota
	Class1xx
	Class2xx
	Class3xx
	Class4xx
	Class5xx
)

func (c Class) String() string {
	switch c {
	case Class1xx:
		return "1xx"
	case Class2xx:
		return "2xx"
	case Class3xx:
		return "3xx"
	case Class4xx:
		return "4xx"
	case Class5xx:
		return "5xx"
	default:
		return "invalid"
	}
}

func Classify(code int) Class {
	switch {
	case code >= 100 && code < 200:
		return Class1xx
	case code >= 200 && code < 300:
		return Class2xx
	case code >= 300 && code < 400:
		return Class3xx
	case code >= 400 && code < 500:
		return Class4xx
	case code >= 500 && code < 600:
		return Class5xx
	default:
		return ClassInvalid
	}
}

func Success(code int) bool { return Classify(code) == Class2xx }

func Redirect(code int) bool { return Classify(code) == Class3xx }

func ClientError(code int) bool { return Classify(code) == Class4xx }

func ServerError(code int) bool { return Classify(code) == Class5xx }

// Retryable 与 retry.RetryableStatus 对齐的判定，供投递层使用。
func Retryable(code int) bool {
	switch code {
	case 408, 425, 429:
		return true
	}
	return ServerError(code)
}
