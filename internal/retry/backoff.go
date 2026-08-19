package retry

import (
	"math"
	"time"
)

// Policy 指数退避策略。
type Policy struct {
	MaxAttempts int
	Base        time.Duration
	Cap         time.Duration
	Factor      float64
	Jitter      float64 // 0..1，乘在 delay 上的随机比例上限；0 表示无抖动
}

func (p Policy) Normalize() Policy {
	if p.MaxAttempts < 1 {
		p.MaxAttempts = 1
	}
	if p.Factor < 1 {
		p.Factor = 2
	}
	if p.Base < 0 {
		p.Base = 0
	}
	if p.Cap < 0 {
		p.Cap = 0
	}
	if p.Jitter < 0 {
		p.Jitter = 0
	}
	if p.Jitter > 1 {
		p.Jitter = 1
	}
	return p
}

// Delay 第 attempt 次失败后的等待（attempt 从 1 起）。最后一次失败不再等待。
func (p Policy) Delay(attempt int) time.Duration {
	p = p.Normalize()
	if attempt < 1 || p.Base == 0 {
		return 0
	}
	exp := float64(attempt - 1)
	mult := math.Pow(p.Factor, exp)
	d := time.Duration(float64(p.Base) * mult)
	if p.Cap > 0 && d > p.Cap {
		d = p.Cap
	}
	if d < 0 {
		return 0
	}
	return d
}

// DelayWithJitter 在 Delay 基础上按 jitter 比例缩小（确定性：用 seed 代替随机，保持可测）。
func (p Policy) DelayWithJitter(attempt int, unit float64) time.Duration {
	d := p.Delay(attempt)
	p = p.Normalize()
	if d == 0 || p.Jitter == 0 {
		return d
	}
	if unit < 0 {
		unit = 0
	}
	if unit > 1 {
		unit = 1
	}
	scale := 1 - p.Jitter*unit
	if scale < 0 {
		scale = 0
	}
	return time.Duration(float64(d) * scale)
}

// ShouldRetry 是否还有下一次尝试。
func (p Policy) ShouldRetry(attempt int) bool {
	p = p.Normalize()
	return attempt < p.MaxAttempts
}
