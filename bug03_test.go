package webhookhub

import (
	"errors"
	"testing"
)

func TestBug03_DispatchAfterCloseNoPanic(t *testing.T) {
	h := New(WithMaxAttempts(1))
	if _, err := h.Subscribe(Endpoint{
		URL:     "http://127.0.0.1:9/hook",
		Events:  []string{"*"},
		Enabled: true,
	}); err != nil {
		t.Fatal(err)
	}
	if err := h.Close(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if rec := recover(); rec != nil {
			t.Fatalf("Dispatch panicked after Close: %v", rec)
		}
	}()
	_, err := h.Dispatch("e", []byte(`{}`))
	if !errors.Is(err, ErrClosed) {
		t.Fatalf("want ErrClosed, got %v", err)
	}
}
