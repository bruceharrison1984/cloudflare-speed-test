package providers

import "fmt"

/* This interface provides methods for retreiving urls necessary for the speed tests */
type UrlProvider struct{}

/* Get the URL for the download speed test */
func (provider *UrlProvider) GetDownloadTestUrl(testId int64, payloadSize int64) string {
	return fmt.Sprintf("https://speed.cloudflare.com/__down?measId=%d&bytes=%d", testId, payloadSize)
}

/* Get the URL for the upload speed test */
func (provider *UrlProvider) GetUploadUrl(testId int64) string {
	return fmt.Sprintf("https://speed.cloudflare.com/__up?measId=%d", testId)
}

/* Get the URL for connection metadata */
func (provider *UrlProvider) GetMetadataUrl() string {
	return "https://speed.cloudflare.com/meta"
}
