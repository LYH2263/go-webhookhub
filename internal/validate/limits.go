package validate

import (
	"strings"
	"unicode"
	"unicode/utf8"

	ierr "github.com/LYH2263/go-webhookhub/internal/errors"
)

const (
	maxURLLen     = 2048
	maxEventLen   = 256
	maxIDLen      = 128
	maxHeaderName = 64
	maxHeaderVal  = 1024
)

func URLLength(raw string) error {
	if len(raw) > maxURLLen {
		return ierr.Wrap(ierr.ErrInvalidURL, "too long")
	}
	return nil
}

func HeaderMap(h map[string]string) error {
	for k, v := range h {
		if strings.TrimSpace(k) == "" {
			return ierr.Wrap(ierr.ErrInvalidEndpoint, "empty header name")
		}
		if len(k) > maxHeaderName || len(v) > maxHeaderVal {
			return ierr.Wrap(ierr.ErrInvalidEndpoint, "header too long")
		}
		if !utf8.ValidString(k) || !utf8.ValidString(v) {
			return ierr.Wrap(ierr.ErrInvalidEndpoint, "header not utf-8")
		}
		for _, r := range k {
			if r < 32 || unicode.IsControl(r) {
				return ierr.Wrap(ierr.ErrInvalidEndpoint, "header control")
			}
		}
	}
	return nil
}

func MaxEventLen(s string) error {
	if len(s) > maxEventLen {
		return ierr.Wrap(ierr.ErrInvalidEvent, "too long")
	}
	return nil
}

func MaxIDLen(id string) error {
	if len(id) > maxIDLen {
		return ierr.Wrap(ierr.ErrInvalidEndpoint, "id too long")
	}
	return nil
}
