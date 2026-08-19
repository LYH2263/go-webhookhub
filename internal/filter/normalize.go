package filter

import "strings"

// Normalize 去掉空白、拆逗号、去重保序。
func Normalize(patterns []string) []string {
	seen := make(map[string]struct{}, len(patterns))
	out := make([]string, 0, len(patterns))
	for _, p := range patterns {
		for _, alt := range strings.Split(p, ",") {
			alt = strings.TrimSpace(alt)
			if alt == "" {
				continue
			}
			if _, ok := seen[alt]; ok {
				continue
			}
			seen[alt] = struct{}{}
			out = append(out, alt)
		}
	}
	return out
}

func IsCatchAll(patterns []string) bool {
	for _, p := range Normalize(patterns) {
		if p == "*" {
			return true
		}
	}
	return false
}
