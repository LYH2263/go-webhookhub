package endpoint

import "github.com/LYH2263/go-webhookhub/internal/filter"

// MatchesEvent 端点是否订阅该事件。禁用端点不匹配。
func MatchesEvent(r Record, event string) bool {
	if !r.Enabled {
		return false
	}
	return filter.Match(r.Events, event)
}

// FilterByEvent 返回匹配且已启用的端点拷贝。
func FilterByEvent(all []Record, event string) []Record {
	var out []Record
	for _, r := range all {
		if MatchesEvent(r, event) {
			out = append(out, CloneRecord(r))
		}
	}
	return out
}
