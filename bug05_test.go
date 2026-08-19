package webhookhub

import (
	"errors"
	"testing"
	"time"
)

func TestBug05_HTTPErrorWrapsSentinel(t *testing.T) {
	h := New(WithMaxAttempts(1), WithFailFast(true), WithHTTPTimeout(2*time.Second))
	defer h.Close()
	if _, err := h.Subscribe(Endpoint{
		URL:     "http://127.0.0.1:1/dead",
		Events:  []string{"*"},
		Enabled: true,
	}); err != nil {
		t.Fatal(err)
	}
	_, err := h.Dispatch("order.created", []byte(`{"ok":1}`))
	if err == nil {
		t.Fatal("expected Dispatch delivery error")
	}
	if !errors.Is(err, ErrHTTP) && !errors.Is(err, ErrTimeout) {
		t.Fatalf("Dispatch err must errors.Is ErrHTTP/ErrTimeout, got %v", err)
	}
}
