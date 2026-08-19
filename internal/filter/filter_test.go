package filter

import "testing"

func TestGlob(t *testing.T) {
	cases := []struct {
		pat, ev string
		ok      bool
	}{
		{"*", "order.created", true},
		{"order.*", "order.created", true},
		{"order.*", "user.created", false},
		{"order.?", "order.x", true},
		{"order.?", "order.xy", false},
		{"exact", "exact", true},
		{"exact", "exact2", false},
	}
	for _, c := range cases {
		if MatchOne(c.pat, c.ev) != c.ok {
			t.Fatalf("%q vs %q want %v", c.pat, c.ev, c.ok)
		}
		comp := Compile(c.pat)
		if comp.Match(c.ev) != c.ok {
			t.Fatalf("compiled %q vs %q want %v", c.pat, c.ev, c.ok)
		}
	}
	if !Match([]string{"a", "b.*"}, "b.1") {
		t.Fatal("list match")
	}
}
