package webhookhub

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBug01_DispatchBodySliceAlias(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	h := New(WithMaxAttempts(1))
	defer h.Close()
	if _, err := h.Subscribe(Endpoint{URL: srv.URL, Events: []string{"*"}, Enabled: true}); err != nil {
		t.Fatal(err)
	}
	body := []byte("hello-payload")
	if _, err := h.Dispatch("ping", body); err != nil {
		t.Fatal(err)
	}
	body[0] = 'X'
	recs := h.log.Recent(1)
	if len(recs) != 1 {
		t.Fatalf("journal len=%d", len(recs))
	}
	if string(recs[0].Payload) != "hello-payload" {
		t.Fatalf("mutating caller body after Dispatch leaked into journal: %q", recs[0].Payload)
	}
}
