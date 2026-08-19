package webhookhub

import ierr "github.com/LYH2263/go-webhookhub/internal/errors"

// 对外哨兵错误。内部包使用同一组变量，便于 errors.Is 跨层判断。
var (
	ErrClosed           = ierr.ErrClosed
	ErrInvalidEndpoint  = ierr.ErrInvalidEndpoint
	ErrInvalidURL       = ierr.ErrInvalidURL
	ErrInvalidEvent     = ierr.ErrInvalidEvent
	ErrNotFound         = ierr.ErrNotFound
	ErrAlreadyExists    = ierr.ErrAlreadyExists
	ErrHTTP             = ierr.ErrHTTP
	ErrStatus           = ierr.ErrStatus
	ErrTimeout          = ierr.ErrTimeout
	ErrCanceled         = ierr.ErrCanceled
	ErrPersist          = ierr.ErrPersist
	ErrNilSigner        = ierr.ErrNilSigner
	ErrEmptyBody        = ierr.ErrEmptyBody
	ErrTooManyEndpoints = ierr.ErrTooManyEndpoints
	ErrDisabled         = ierr.ErrDisabled
)
