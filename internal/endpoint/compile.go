package endpoint

import (
	"strings"

	"github.com/LYH2263/go-webhookhub/internal/filter"
)

// Compiled 带预编译过滤器的端点。
type Compiled struct {
	Record Record
	Filt   []filter.Compiled
}

func CompileRecord(r Record) Compiled {
	r = CloneRecord(r)
	return Compiled{Record: r, Filt: filter.CompileAll(r.Events)}
}

func (c Compiled) Matches(event string) bool {
	if !c.Record.Enabled {
		return false
	}
	if len(c.Filt) == 0 {
		return filter.MatchOne("*", event)
	}
	return filter.MatchCompiled(c.Filt, event)
}

func CompileAll(recs []Record) []Compiled {
	out := make([]Compiled, 0, len(recs))
	for _, r := range recs {
		out = append(out, CompileRecord(r))
	}
	return out
}

func SplitEvents(s string) []string {
	parts := strings.FieldsFunc(s, func(r rune) bool {
		return r == ',' || r == ' ' || r == ';'
	})
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
