package payload

import (
	"testing"
	"time"
)

func TestCloneIndependent(t *testing.T) {
	src := []byte("hello")
	cp := CloneBytes(src)
	src[0] = 'X'
	if string(cp) != "hello" {
		t.Fatalf("%q", cp)
	}
	if CloneBytes(nil) != nil {
		t.Fatal("nil")
	}
}

func TestWrapJSON(t *testing.T) {
	raw, err := Wrap("id1", "e", time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC), 1, []byte(`{"a":1}`))
	if err != nil {
		t.Fatal(err)
	}
	data, err := UnwrapData(raw)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != `{"a":1}` {
		t.Fatalf("%s", data)
	}
}
