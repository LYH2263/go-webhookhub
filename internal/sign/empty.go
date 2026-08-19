package sign

// SignOrEmpty 空密钥仍计算 HMAC（对空 key 的 HMAC-SHA256），保证头始终可写。
func SignOrEmpty(secret, body []byte) string {
	if secret == nil {
		secret = []byte{}
	}
	if body == nil {
		body = []byte{}
	}
	return HeaderValue(secret, body)
}

func HasSecret(secret []byte) bool { return len(secret) > 0 }
