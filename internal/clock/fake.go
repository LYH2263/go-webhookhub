package clock

import (
	"sync"
	"time"
)

// Fake 可推进的测试时钟。
type Fake struct {
	mu sync.Mutex
	t  time.Time
}

func NewFake(t time.Time) *Fake {
	if t.IsZero() {
		t = time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	}
	return &Fake{t: t.UTC()}
}

func (f *Fake) Now() time.Time {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.t
}

func (f *Fake) Set(t time.Time) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.t = t.UTC()
}

func (f *Fake) Advance(d time.Duration) time.Time {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.t = f.t.Add(d)
	return f.t
}

func (f *Fake) Unix() int64 {
	return f.Now().Unix()
}
