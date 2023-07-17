package clients

import (
	"net/http"
)

type SpeedTestTransport struct {
	base http.RoundTripper
	http.Transport
}

func (t *SpeedTestTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.MaxIdleConns = 100
	t.MaxIdleConnsPerHost = 100

	req.Header.Add("User-Agent", "go-cf-speed-test")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Content-Type", "text/plain;charset=UTF-8")

	if t.base == nil {
		t.base = http.DefaultTransport
	}

	return t.base.RoundTrip(req)
}
