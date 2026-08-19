package filter

import "strings"

// Match 判断 event 是否命中任一模式。空模式列表视为不匹配。
func Match(patterns []string, event string) bool {
	if event == "" {
		return false
	}
	for _, p := range patterns {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		for _, alt := range strings.Split(p, ",") {
			alt = strings.TrimSpace(alt)
			if alt == "" {
				continue
			}
			if MatchOne(alt, event) {
				return true
			}
		}
	}
	return false
}

// MatchOne 单模式：* 任意串（含点），? 单字符，其余字面。
func MatchOne(pattern, event string) bool {
	if pattern == "*" {
		return event != ""
	}
	return matchAt(pattern, event)
}

func matchAt(p, s string) bool {
	for len(p) > 0 {
		switch p[0] {
		case '*':
			p = p[1:]
			if p == "" {
				return true
			}
			for i := 0; i <= len(s); i++ {
				if matchAt(p, s[i:]) {
					return true
				}
			}
			return false
		case '?':
			if len(s) == 0 {
				return false
			}
			p = p[1:]
			s = s[1:]
		default:
			if len(s) == 0 || p[0] != s[0] {
				return false
			}
			p = p[1:]
			s = s[1:]
		}
	}
	return len(s) == 0
}

// AnyStar 是否包含通配。
func AnyStar(patterns []string) bool {
	for _, p := range patterns {
		if strings.ContainsAny(p, "*?") {
			return true
		}
	}
	return false
}
