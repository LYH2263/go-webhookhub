package api

import (
	"context"
	"net"
	"net/http"
	"time"
)

// ListenAndServe 在 addr 上启动，可配合 Shutdown。
func ListenAndServe(ctx context.Context, addr string, h http.Handler) error {
	srv := &http.Server{
		Addr:              addr,
		Handler:           h,
		ReadHeaderTimeout: 5 * time.Second,
	}
	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.ListenAndServe()
	}()
	select {
	case <-ctx.Done():
		c2, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		_ = srv.Shutdown(c2)
		return ctx.Err()
	case err := <-errCh:
		if err == http.ErrServerClosed {
			return nil
		}
		return err
	}
}

func LocalAddr(ln net.Listener) string {
	if ln == nil {
		return ""
	}
	return ln.Addr().String()
}
