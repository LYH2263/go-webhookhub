package webhookhub

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestBug10_CloseFlushesJournalBeforeRelease(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "journal.jsonl")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	h := New(WithJournalPath(path), WithMaxAttempts(1))
	if _, err := h.Subscribe(Endpoint{URL: srv.URL, Events: []string{"*"}, Enabled: true}); err != nil {
		t.Fatal(err)
	}
	if _, err := h.Dispatch("order.created", []byte(`{}`)); err != nil {
		t.Fatal(err)
	}
	if err := h.Close(); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(raw, []byte("order.created")) {
		t.Fatalf("journal file missing records after Close (not flushed): %q", raw)
	}
}
