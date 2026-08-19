package httpx

import "testing"

func TestClassify(t *testing.T) {
	if !Success(201) || Retryable(503) == false || Retryable(400) {
		t.Fatal("class")
	}
	if Classify(404).String() != "4xx" {
		t.Fatal(Classify(404))
	}
}
