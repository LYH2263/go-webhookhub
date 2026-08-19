package clock

import "time"

// Clock 可注入时钟，便于日志时间戳与测试冻结。
type Clock interface {
	Now() time.Time
}

// AsUTC 把时间规范到 UTC；零值保持零值。
func AsUTC(t time.Time) time.Time {
	if t.IsZero() {
		return t
	}
	return t.UTC()
}

// UnixMilli 毫秒时间戳。
func UnixMilli(c Clock) int64 {
	if c == nil {
		return time.Now().UTC().UnixMilli()
	}
	return c.Now().UTC().UnixMilli()
}
