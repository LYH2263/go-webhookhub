package clock

import "time"

// Real 使用 time.Now().UTC()。
type Real struct{}

func (Real) Now() time.Time { return time.Now().UTC() }

// RealMS 与 Real 相同，单独类型便于选项区分。
type RealMS struct{}

func (RealMS) Now() time.Time { return time.Now().UTC().Truncate(time.Millisecond) }
