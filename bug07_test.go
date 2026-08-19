package webhookhub

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"
)

type bug07RT struct {
	n      *int
	cancel context.CancelFunc
}

func (rt *bug07RT) RoundTrip(req *http.Request) (*http.Response, error) {
	*rt.n++
	if *rt.n == 1 && rt.cancel != nil {
		rt.cancel()
	}
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(nil)),
		Header:     make(http.Header),
		Request:    req,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
	}, nil
}

func TestBug07_DispatchContextHonorsCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	n := 0
	h := New(WithTransport(&bug07RT{n: &n, cancel: cancel}), WithMaxAttempts(1))
	defer h.Close()
	for i := 0; i < 2; i++ {
		if _, err := h.Subscribe(Endpoint{
			URL:     "http://127.0.0.1:9/hook",
			Events:  []string{"*"},
			Enabled: true,
		}); err != nil {
			t.Fatal(err)
		}
	}
	_, err := h.DispatchContext(ctx, "e", []byte(`{}`))
	if n > 1 {
		t.Fatalf("fanout ignored ctx cancel: deliveries=%d", n)
	}
	if err == nil {
		t.Fatal("expected context cancel error")
	}
}
