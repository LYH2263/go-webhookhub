package webhookhub

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/LYH2263/go-webhookhub/internal/sign"
)

func TestDispatchHappyPath(t *testing.T) {
	var gotSig, gotEvent string
	var gotBody []byte
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotSig = r.Header.Get(HeaderSignature)
		gotEvent = r.Header.Get(HeaderEvent)
		gotBody, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	h := New(WithMaxAttempts(1), WithHTTPTimeout(2*time.Second))
	defer h.Close()

	secret := []byte("topsecret")
	id, err := h.Subscribe(Endpoint{
		URL:     srv.URL,
		Secret:  secret,
		Events:  []string{"order.*"},
		Enabled: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if id == "" {
		t.Fatal("empty id")
	}

	body := []byte(`{"n":1}`)
	results, err := h.Dispatch("order.created", body)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 || !results[0].OK {
		t.Fatalf("results=%+v", results)
	}
	if gotEvent != "order.created" {
		t.Fatalf("event header %q", gotEvent)
	}
	if !sign.Verify([]byte("topsecret"), gotBody, gotSig) {
		t.Fatalf("bad signature %q body %s", gotSig, gotBody)
	}
	recs := h.RecentDeliveries(10)
	if len(recs) != 1 || !recs[0].OK {
		t.Fatalf("journal %+v", recs)
	}
}

func TestPayloadAndSecretCloned(t *testing.T) {
	var got []byte
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got, _ = io.ReadAll(r.Body)
		w.WriteHeader(200)
	}))
	defer srv.Close()

	h := New(WithMaxAttempts(1))
	defer h.Close()
	secret := []byte("abcdef")
	_, err := h.Subscribe(Endpoint{URL: srv.URL, Secret: secret, Events: []string{"*"}, Enabled: true})
	if err != nil {
		t.Fatal(err)
	}
	secret[0] = 'Z'
	views := h.ListEndpoints()
	if string(views[0].Secret) != "abcdef" {
		t.Fatalf("secret alias: %q", views[0].Secret)
	}
	views[0].Secret[0] = 'Y'
	views2 := h.ListEndpoints()
	if string(views2[0].Secret) != "abcdef" {
		t.Fatalf("view alias: %q", views2[0].Secret)
	}

	body := []byte(`{"k":"v"}`)
	if _, err := h.Dispatch("ping", body); err != nil {
		t.Fatal(err)
	}
	body[2] = 'X'
	var env map[string]json.RawMessage
	if err := json.Unmarshal(got, &env); err != nil {
		t.Fatal(err)
	}
	if string(env["data"]) != `{"k":"v"}` {
		t.Fatalf("payload alias through envelope: %s", env["data"])
	}
}

func TestCloseThenDispatch(t *testing.T) {
	h := New()
	if err := h.Close(); err != nil {
		t.Fatal(err)
	}
	_, err := h.Dispatch("x", []byte(`{}`))
	if !errors.Is(err, ErrClosed) {
		t.Fatalf("want ErrClosed, got %v", err)
	}
	_, err = h.Subscribe(Endpoint{URL: "http://127.0.0.1:1/x", Events: []string{"*"}, Enabled: true})
	if !errors.Is(err, ErrClosed) {
		t.Fatalf("subscribe want ErrClosed, got %v", err)
	}
}

func TestDispatchContextCanceled(t *testing.T) {
	h := New(WithMaxAttempts(1))
	defer h.Close()
	_, err := h.Subscribe(Endpoint{
		URL:     "http://127.0.0.1:1/nowhere",
		Events:  []string{"*"},
		Enabled: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err = h.DispatchContext(ctx, "e", []byte(`{}`))
	if err == nil {
		t.Fatal("expected cancel error")
	}
}

func TestEventFilterNoMatch(t *testing.T) {
	h := New()
	defer h.Close()
	_, err := h.Subscribe(Endpoint{
		URL:     "http://127.0.0.1:9/x",
		Events:  []string{"order.created"},
		Enabled: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	res, err := h.Dispatch("user.signup", []byte(`{}`))
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 0 {
		t.Fatalf("want 0 results, got %+v", res)
	}
}

func TestSubscribeRollbackOnPersistFail(t *testing.T) {
	sink := &failSink{}
	h := New(WithEndpointSink(sink))
	defer h.Close()
	_, err := h.Subscribe(Endpoint{
		URL:     "http://127.0.0.1:9/x",
		Events:  []string{"*"},
		Enabled: true,
	})
	if err == nil {
		t.Fatal("expected persist error")
	}
	if h.EndpointCount() != 0 {
		t.Fatalf("did not rollback, n=%d", h.EndpointCount())
	}
}

type failSink struct{}

func (*failSink) SaveEndpoints([]EndpointView) error {
	return errors.New("disk full")
}

func TestDefaultSignerInstalled(t *testing.T) {
	h := New(WithSigner(nil))
	defer h.Close()
	if h.signer == nil {
		t.Fatal("signer is nil")
	}
}

func TestDisabledEndpointSkipped(t *testing.T) {
	called := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called++
		w.WriteHeader(200)
	}))
	defer srv.Close()
	h := New(WithMaxAttempts(1))
	defer h.Close()
	id, err := h.Subscribe(Endpoint{URL: srv.URL, Events: []string{"*"}, Enabled: false})
	if err != nil {
		t.Fatal(err)
	}
	res, err := h.Dispatch("e", []byte(`{}`))
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 0 || called != 0 {
		t.Fatalf("disabled still delivered: %+v called=%d", res, called)
	}
	if err := h.SetEnabled(id, true); err != nil {
		t.Fatal(err)
	}
	res, err = h.Dispatch("e", []byte(`{}`))
	if err != nil || len(res) != 1 || !res[0].OK {
		t.Fatalf("%v %+v", err, res)
	}
}

func TestHTTPErrorWrapped(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()
	h := New(WithMaxAttempts(1), WithBackoff(0, 0))
	defer h.Close()
	_, err := h.Subscribe(Endpoint{URL: srv.URL, Events: []string{"*"}, Enabled: true})
	if err != nil {
		t.Fatal(err)
	}
	res, err := h.Dispatch("e", []byte(`{}`))
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 1 || res[0].OK {
		t.Fatalf("%+v", res)
	}
	if !errors.Is(errors.New(res[0].Err), ErrStatus) && res[0].StatusCode != 500 {
		if res[0].StatusCode != 500 {
			t.Fatalf("status %d err %s", res[0].StatusCode, res[0].Err)
		}
	}
}

func TestUnsubscribe(t *testing.T) {
	h := New()
	defer h.Close()
	id, err := h.Subscribe(Endpoint{URL: "http://127.0.0.1:9/x", Events: []string{"*"}, Enabled: true})
	if err != nil {
		t.Fatal(err)
	}
	if err := h.Unsubscribe(id); err != nil {
		t.Fatal(err)
	}
	if h.EndpointCount() != 0 {
		t.Fatal("still registered")
	}
}

func TestStats(t *testing.T) {
	h := New()
	defer h.Close()
	st := h.Stats()
	if st.Closed || st.DefaultMaxTry < 1 {
		t.Fatalf("%+v", st)
	}
}

func TestDispatchTimeoutOption(t *testing.T) {
	h := New(WithHTTPTimeout(50 * time.Millisecond), WithMaxAttempts(1))
	defer h.Close()
	if h.httpTimeout != 50*time.Millisecond {
		t.Fatal(h.httpTimeout)
	}
}
