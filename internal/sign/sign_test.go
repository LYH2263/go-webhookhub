package sign

import "testing"

func TestHMACRoundTrip(t *testing.T) {
	secret := []byte("k")
	body := []byte(`{"a":1}`)
	hv := HeaderValue(secret, body)
	if !Verify(secret, body, hv) {
		t.Fatal("verify failed")
	}
	if Verify([]byte("other"), body, hv) {
		t.Fatal("wrong secret accepted")
	}
	s := Default()
	if s.HeaderName() != HeaderName {
		t.Fatal(s.HeaderName())
	}
	if s.HeaderValue(secret, body) != hv {
		t.Fatal("mismatch")
	}
}
