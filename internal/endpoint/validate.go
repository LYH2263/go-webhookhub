package endpoint

import (
	"strings"

	"github.com/LYH2263/go-webhookhub/internal/idgen"
	"github.com/LYH2263/go-webhookhub/internal/validate"
)

func NormalizeID(id string) string {
	id = strings.TrimSpace(id)
	if id == "" {
		return idgen.EndpointID()
	}
	return id
}

func NormalizeEvents(ev []string) []string {
	if len(ev) == 0 {
		return []string{"*"}
	}
	seen := make(map[string]struct{}, len(ev))
	out := make([]string, 0, len(ev))
	for _, e := range ev {
		e = strings.TrimSpace(e)
		if e == "" {
			continue
		}
		if _, ok := seen[e]; ok {
			continue
		}
		seen[e] = struct{}{}
		out = append(out, e)
	}
	if len(out) == 0 {
		return []string{"*"}
	}
	return out
}

func Validate(r Record, allowHTTP bool) error {
	if err := validate.URL(r.URL, allowHTTP); err != nil {
		return err
	}
	if r.MaxAttempts < 0 {
		r.MaxAttempts = 0
	}
	for _, e := range r.Events {
		if strings.TrimSpace(e) == "" {
			continue
		}
		if err := validate.EventName(strings.TrimSpace(e)); err != nil {
			return err
		}
	}
	return nil
}
