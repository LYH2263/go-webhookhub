package webhookhub

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBug04_NilSignerNoPanic(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	h := New(WithMaxAttempts(1), WithSigner(nil))
	defer h.Close()
	if _, err := h.Subscribe(Endpoint{URL: srv.URL, Events: []string{"*"}, Enabled: true}); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if rec := recover(); rec != nil {
			t.Fatalf("Dispatch panicked on nil signer: %v", rec)
		}
	}()
	res, err := h.Dispatch("e", []byte(`{}`))
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 1 || !res[0].OK {
		t.Fatalf("results=%+v", res)
	}
	if !res[0].Signed {
		t.Fatal("expected default HMAC signer to produce a signature")
	}
}
