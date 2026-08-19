package retry

import (
	"context"
	"testing"
	"time"
)

func TestBackoffCap(t *testing.T) {
	p := Policy{Base: 10 * time.Millisecond, Cap: 50 * time.Millisecond, Factor: 2, MaxAttempts: 8}.Normalize()
	d := p.Delay(8)
	if d > 50*time.Millisecond {
		t.Fatalf("cap failed: %s", d)
	}
	if p.Delay(1) != 10*time.Millisecond {
		t.Fatal(p.Delay(1))
	}
}

func TestWaitCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := Wait(ctx, time.Second); err == nil {
		t.Fatal("expected cancel")
	}
}

func TestWaitZero(t *testing.T) {
	if err := Wait(context.Background(), 0); err != nil {
		t.Fatal(err)
	}
}
