package webhookhub

import "testing"

func TestBug06_SubscribePersistFailureRollsBack(t *testing.T) {
	sink := &failSink{}
	h := New(WithEndpointSink(sink))
	defer h.Close()
	id, err := h.Subscribe(Endpoint{
		URL:     "http://127.0.0.1:9/x",
		Events:  []string{"*"},
		Enabled: true,
	})
	if err == nil {
		t.Fatal("expected persist error, Subscribe returned success")
	}
	if id != "" {
		t.Fatalf("id=%q on persist failure", id)
	}
	if h.EndpointCount() != 0 {
		t.Fatalf("partial subscribe left in registry, n=%d", h.EndpointCount())
	}
}
