package journal

import "testing"

func TestRingRecent(t *testing.T) {
	l := New(4, nil, "")
	for i := 0; i < 6; i++ {
		l.Append(Record{ID: string(rune('a' + i)), Event: "e"})
	}
	got := l.Recent(3)
	if len(got) != 3 {
		t.Fatalf("len %d", len(got))
	}
	if l.Len() != 4 {
		t.Fatalf("cap ring len %d", l.Len())
	}
	if err := l.Flush(); err != nil {
		t.Fatal(err)
	}
	if err := l.Sync(); err != nil {
		t.Fatal(err)
	}
	if err := l.Close(); err != nil {
		t.Fatal(err)
	}
}
