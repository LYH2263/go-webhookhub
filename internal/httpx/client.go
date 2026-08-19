package httpx

import (
	"net/http"
	"time"
)

// Client 带超时与 UA 的 HTTP 客户端包装。Hub.Close 会将其置 nil。
type Client struct {
	hc        *http.Client
	userAgent string
}

func New(timeout time.Duration, rt http.RoundTripper, ua string) *Client {
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	if rt == nil {
		rt = http.DefaultTransport
	}
	return &Client{
		hc: &http.Client{
			Timeout:   timeout,
			Transport: rt,
		},
		userAgent: ua,
	}
}

func (c *Client) HTTP() *http.Client {
	if c == nil {
		return nil
	}
	return c.hc
}

func (c *Client) UserAgent() string {
	if c == nil {
		return ""
	}
	return c.userAgent
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	if req != nil && c.userAgent != "" && req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", c.userAgent)
	}
	return c.hc.Do(req)
}

func (c *Client) CloseIdle() {
	if c == nil || c.hc == nil {
		return
	}
	c.hc.CloseIdleConnections()
}

func (c *Client) Timeout() time.Duration {
	if c == nil || c.hc == nil {
		return 0
	}
	return c.hc.Timeout
}

func errNil() error {
	return errNilClient
}

var errNilClient = &nilClientError{}

type nilClientError struct{}

func (*nilClientError) Error() string { return "webhookhub: nil http client" }
