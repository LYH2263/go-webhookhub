package endpoint

import "testing"

func TestRegistryClone(t *testing.T) {
	r := NewRegistry(10, true)
	id, err := r.Add(Record{URL: "http://127.0.0.1:9/h", Secret: []byte("abc"), Events: []string{"x.*"}, Enabled: true})
	if err != nil {
		t.Fatal(err)
	}
	got, err := r.Get(id)
	if err != nil {
		t.Fatal(err)
	}
	got.Secret[0] = 'Z'
	got2, _ := r.Get(id)
	if string(got2.Secret) != "abc" {
		t.Fatalf("%q", got2.Secret)
	}
	m := r.Match("x.1")
	if len(m) != 1 {
		t.Fatalf("%d", len(m))
	}
	if len(r.Match("y.1")) != 0 {
		t.Fatal("false match")
	}
}
