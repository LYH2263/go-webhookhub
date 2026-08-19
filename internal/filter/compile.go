package filter

import "strings"

type kind int

const (
	kindLit kind = iota
	kindAny
	kindOne
)

type token struct {
	k kind
	s string
}

// Compiled 预编译的事件模式，避免每次扇出重复切分。
type Compiled struct {
	raw    string
	tokens []token
}

func Compile(pattern string) Compiled {
	pattern = strings.TrimSpace(pattern)
	var toks []token
	var lit strings.Builder
	flush := func() {
		if lit.Len() == 0 {
			return
		}
		toks = append(toks, token{k: kindLit, s: lit.String()})
		lit.Reset()
	}
	for i := 0; i < len(pattern); i++ {
		switch pattern[i] {
		case '*':
			flush()
			toks = append(toks, token{k: kindAny})
		case '?':
			flush()
			toks = append(toks, token{k: kindOne})
		default:
			lit.WriteByte(pattern[i])
		}
	}
	flush()
	return Compiled{raw: pattern, tokens: toks}
}

func CompileAll(patterns []string) []Compiled {
	var out []Compiled
	for _, p := range patterns {
		for _, alt := range strings.Split(p, ",") {
			alt = strings.TrimSpace(alt)
			if alt == "" {
				continue
			}
			out = append(out, Compile(alt))
		}
	}
	return out
}

func (c Compiled) Match(event string) bool {
	if c.raw == "*" {
		return event != ""
	}
	return matchTokens(c.tokens, event)
}

func matchTokens(toks []token, s string) bool {
	if len(toks) == 0 {
		return s == ""
	}
	t := toks[0]
	rest := toks[1:]
	switch t.k {
	case kindAny:
		for i := 0; i <= len(s); i++ {
			if matchTokens(rest, s[i:]) {
				return true
			}
		}
		return false
	case kindOne:
		if len(s) == 0 {
			return false
		}
		return matchTokens(rest, s[1:])
	default:
		if !strings.HasPrefix(s, t.s) {
			return false
		}
		return matchTokens(rest, s[len(t.s):])
	}
}

func MatchCompiled(cs []Compiled, event string) bool {
	for _, c := range cs {
		if c.Match(event) {
			return true
		}
	}
	return false
}

func (c Compiled) Raw() string { return c.raw }
