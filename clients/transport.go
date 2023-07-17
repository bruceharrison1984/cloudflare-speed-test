package clients

import (
	"net/http"
)

/* Overrides for the default transport */
type cloudflareSpeedTestTransport struct {
	base http.RoundTripper
	http.Transport
}

func NewCloudflareSpeedTestTransport() *cloudflareSpeedTestTransport {
	return &cloudflareSpeedTestTransport{}
}

func (t *cloudflareSpeedTestTransport) RoundTrip(req *http.Request) (*http.Response, error) {
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
