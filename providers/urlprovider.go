package providers

import "fmt"

type UrlProvider struct{}

func (provider *UrlProvider) GetDownloadTestUrl(testId int64, payloadSize int64) string {
	return fmt.Sprintf("https://speed.cloudflare.com/__down?measId=%d&bytes=%d", testId, payloadSize)
}

func (provider *UrlProvider) GetUploadUrl(testId int64) string {
	return fmt.Sprintf("https://speed.cloudflare.com/__up?measId=%d", testId)
}

func (provider *UrlProvider) GetMetadataUrl() string {
	return "https://speed.cloudflare.com/meta"
}
