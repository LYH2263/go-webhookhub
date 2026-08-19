package deliver

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/LYH2263/go-webhookhub/internal/httpx"
	"github.com/LYH2263/go-webhookhub/internal/sign"
)

type closeCounter struct {
	n int
}

type countingBody struct {
	io.ReadCloser
	c *closeCounter
}

func (b *countingBody) Close() error {
	b.c.n++
	return b.ReadCloser.Close()
}

type countingRT struct {
	c *closeCounter
}

func (rt *countingRT) RoundTrip(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       &countingBody{ReadCloser: io.NopCloser(bytes.NewReader(nil)), c: rt.c},
		Header:     make(http.Header),
		Request:    req,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
	}, nil
}

func TestBug09_ResponseBodyClosed(t *testing.T) {
	ctr := &closeCounter{}
	client := httpx.New(2*time.Second, &countingRT{c: ctr}, "test")
	p := &Poster{Client: client, Signer: sign.Default()}
	res := p.Post(context.Background(), Request{
		URL:        "http://127.0.0.1/hook",
		Body:       []byte(`{}`),
		Event:      "e",
		DeliveryID: "d1",
		Timeout:    time.Second,
		Attempt:    1,
		EndpointID: "ep1",
	})
	if !res.OK {
		t.Fatalf("post %+v", res)
	}
	if ctr.n < 1 {
		t.Fatalf("response body Close count=%d, want >=1", ctr.n)
	}
}
