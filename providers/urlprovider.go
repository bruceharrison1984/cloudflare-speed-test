package providers

import "fmt"

const (
	DOWNLOAD_TEST_URL = "https://speed.cloudflare.com/__down?measId=%d&bytes=%d"
	UPLOAD_TEST_URL   = "https://speed.cloudflare.com/__up?measId=%d"
	METADATA_URL      = "https://speed.cloudflare.com/meta"
)

/* This interface provides methods for retreiving urls necessary for the speed tests */
type UrlProvider struct{}

/* Get the URL for the download speed test */
func (provider *UrlProvider) GetDownloadTestUrl(testId int64, payloadSize int64) string {
	return fmt.Sprintf(DOWNLOAD_TEST_URL, testId, payloadSize)
}

/* Get the URL for the upload speed test */
func (provider *UrlProvider) GetUploadUrl(testId int64) string {
	return fmt.Sprintf(UPLOAD_TEST_URL, testId)
}

/* Get the URL for connection metadata */
func (provider *UrlProvider) GetMetadataUrl() string {
	return METADATA_URL
}
