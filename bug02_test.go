package webhookhub

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LYH2263/go-webhookhub/internal/sign"
)

func TestBug02_ListEndpointsSecretSliceAlias(t *testing.T) {
	var gotSig string
	var gotBody []byte
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotSig = r.Header.Get(HeaderSignature)
		gotBody, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	h := New(WithMaxAttempts(1))
	defer h.Close()
	secret := []byte("s3cret-key!!")
	if _, err := h.Subscribe(Endpoint{
		URL:     srv.URL,
		Secret:  append([]byte(nil), secret...),
		Events:  []string{"*"},
		Enabled: true,
	}); err != nil {
		t.Fatal(err)
	}
	views := h.ListEndpoints()
	if len(views) != 1 || len(views[0].Secret) == 0 {
		t.Fatalf("views=%+v", views)
	}
	views[0].Secret[0] ^= 0xFF
	if _, err := h.Dispatch("ping", []byte(`{"a":1}`)); err != nil {
		t.Fatal(err)
	}
	if !sign.Verify(secret, gotBody, gotSig) {
		t.Fatalf("mutating ListEndpoints secret changed later signing: sig=%q", gotSig)
	}
}
