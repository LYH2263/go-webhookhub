package retry

import (
	"context"
	"errors"
	"testing"
)

func TestBug08_RetryWaitHonorsContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	calls := 0
	_, _, err := Loop(ctx, Policy{MaxAttempts: 5, Base: 0, Factor: 2}, nil, func(attempt int) (bool, int, error) {
		calls++
		return false, 503, errors.New("down")
	})
	if calls != 0 {
		t.Fatalf("Loop ignored canceled ctx, fn calls=%d", calls)
	}
	if err == nil {
		t.Fatal("Loop expected cancel error")
	}
	if err := Wait(ctx, 0); err == nil {
		t.Fatal("Wait ignored already-canceled ctx")
	}
}
