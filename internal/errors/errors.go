package errors

import (
	"errors"
	"fmt"
)

var (
	ErrClosed           = errors.New("webhookhub: closed")
	ErrInvalidEndpoint  = errors.New("webhookhub: invalid endpoint")
	ErrInvalidURL       = errors.New("webhookhub: invalid url")
	ErrInvalidEvent     = errors.New("webhookhub: invalid event")
	ErrNotFound         = errors.New("webhookhub: not found")
	ErrAlreadyExists    = errors.New("webhookhub: already exists")
	ErrHTTP             = errors.New("webhookhub: http delivery")
	ErrStatus           = errors.New("webhookhub: unexpected status")
	ErrTimeout          = errors.New("webhookhub: timeout")
	ErrCanceled         = errors.New("webhookhub: canceled")
	ErrPersist          = errors.New("webhookhub: persist")
	ErrNilSigner        = errors.New("webhookhub: nil signer")
	ErrEmptyBody        = errors.New("webhookhub: empty body")
	ErrTooManyEndpoints = errors.New("webhookhub: too many endpoints")
	ErrDisabled         = errors.New("webhookhub: endpoint disabled")
	ErrNilClient        = errors.New("webhookhub: nil http client")
	ErrNilJournal       = errors.New("webhookhub: nil journal")
	ErrFlush            = errors.New("webhookhub: journal flush")
	ErrSync             = errors.New("webhookhub: journal sync")
)

// Wrap 把哨兵与说明拼成可 errors.Is 的错误。
func Wrap(sentinel error, msg string) error {
	if sentinel == nil {
		return errors.New(msg)
	}
	if msg == "" {
		return sentinel
	}
	return fmt.Errorf("%w: %s", sentinel, msg)
}

// WrapErr 用 %w 同时包裹哨兵与底层错误，errors.Is 对两者都成立。
func WrapErr(sentinel, err error) error {
	if err == nil {
		return sentinel
	}
	if sentinel == nil {
		return err
	}
	return fmt.Errorf("%w: %w", sentinel, err)
}

// Wrapf 带格式化说明的哨兵包装。
func Wrapf(sentinel error, format string, args ...any) error {
	msg := fmt.Sprintf(format, args...)
	return Wrap(sentinel, msg)
}

// IsAny 判断 err 是否匹配任一哨兵。
func IsAny(err error, sentinels ...error) bool {
	for _, s := range sentinels {
		if errors.Is(err, s) {
			return true
		}
	}
	return false
}
