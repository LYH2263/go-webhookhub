package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/LYH2263/go-webhookhub"
	ierr "github.com/LYH2263/go-webhookhub/internal/errors"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	_ = enc.Encode(v)
}

type errBody struct {
	Error string `json:"error"`
}

func writeErr(w http.ResponseWriter, err error) {
	status := http.StatusBadRequest
	switch {
	case errors.Is(err, webhookhub.ErrNotFound), errors.Is(err, ierr.ErrNotFound):
		status = http.StatusNotFound
	case errors.Is(err, webhookhub.ErrClosed), errors.Is(err, ierr.ErrClosed):
		status = http.StatusServiceUnavailable
	case errors.Is(err, webhookhub.ErrAlreadyExists), errors.Is(err, ierr.ErrAlreadyExists):
		status = http.StatusConflict
	case errors.Is(err, webhookhub.ErrCanceled), errors.Is(err, ierr.ErrCanceled):
		status = http.StatusRequestTimeout
	case errors.Is(err, webhookhub.ErrTimeout), errors.Is(err, ierr.ErrTimeout):
		status = http.StatusGatewayTimeout
	case err == nil:
		status = http.StatusOK
	}
	msg := "error"
	if err != nil {
		msg = err.Error()
	}
	writeJSON(w, status, errBody{Error: msg})
}

func readJSON(r *http.Request, v any) error {
	defer r.Body.Close()
	dec := json.NewDecoder(io.LimitReader(r.Body, 1<<20))
	dec.DisallowUnknownFields()
	if err := dec.Decode(v); err != nil {
		return ierr.Wrap(ierr.ErrInvalidEndpoint, err.Error())
	}
	return nil
}
