package validate

import (
	"net/url"
	"strings"

	ierr "github.com/LYH2263/go-webhookhub/internal/errors"
)

// URL 校验投递地址。allowHTTP 为 false 时仅允许 https。
func URL(raw string, allowHTTP bool) error {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ierr.Wrap(ierr.ErrInvalidURL, "empty")
	}
	u, err := url.Parse(raw)
	if err != nil {
		return ierr.WrapErr(ierr.ErrInvalidURL, err)
	}
	switch strings.ToLower(u.Scheme) {
	case "https":
	case "http":
		if !allowHTTP {
			return ierr.Wrap(ierr.ErrInvalidURL, "http not allowed")
		}
	default:
		return ierr.Wrap(ierr.ErrInvalidURL, "scheme "+u.Scheme)
	}
	if u.Host == "" {
		return ierr.Wrap(ierr.ErrInvalidURL, "missing host")
	}
	if strings.Contains(u.Host, " ") {
		return ierr.Wrap(ierr.ErrInvalidURL, "invalid host")
	}
	return nil
}

func EventName(s string) error {
	s = strings.TrimSpace(s)
	if s == "" {
		return ierr.ErrInvalidEvent
	}
	for _, c := range s {
		if c < 32 || c == 127 {
			return ierr.Wrap(ierr.ErrInvalidEvent, "control char")
		}
	}
	return nil
}

func NonEmptyID(id string) error {
	if strings.TrimSpace(id) == "" {
		return ierr.Wrap(ierr.ErrInvalidEndpoint, "empty id")
	}
	return nil
}
