package types

import "time"

type BandwidthTestType int16

const (
	Upload   BandwidthTestType = 0
	Download BandwidthTestType = 1
)

/** This is the calculated results of a speed test, based on the RawBandwidthClientResult */
type BandwidthClientResult struct {
	ResultType       BandwidthTestType
	PayloadSizeBytes int64
	ServerTiming     float64 // timing value in seconds
	Ping             float64 // timing value in seconds
	TransferDuration float64 // timing value in seconds
	SpeedMbps        float64 // calculated speed. Upload/Download determined by result type
	Ttfb             float64 // timing value in seconds
	Ttlb             float64 // timing value in seconds
}

/** This is the raw results directly produced from a speed test web request */
type RawBandwidthClientResult struct {
	ResultType       BandwidthTestType
	TTFB             time.Duration
	TTLB             time.Duration
	ServerTiming     time.Duration
	PayloadSizeBytes int64
	LoggingToken     string
}

/** These are parameters for running an individual speed test run */
type SpeedTestCase struct {
	PayloadSize int64
	Iterations  int
	TestType    BandwidthTestType
}

/** Final results summary from a full speed test run. Based on percentiles of BandwidthClientResult[] */
type BandwidthClientResultSummary struct {
	Ping              float64 `json:"ping"`
	DownloadSpeedMbps float64 `json:"downloadSpeedMbps"`
	UploadSpeedMbps   float64 `json:"uploadSpeedMbps"`
}

type SpeedTestSummary struct {
	Metadata    *CloudflareMetadata           `json:"metadata"`
	Bandwidth   *BandwidthClientResultSummary `json:"bandwidth"`
	TestResults []*BandwidthClientResult      `json:"testResults"`
}

// The is the response data returned from the Cloudflare metadata endpoint
type CloudflareMetadata struct {
	Hostname       string `json:"hostname"`
	ClientIp       string `json:"clientIp"`
	HttpProtocol   string `json:"httpProtocol"`
	Asn            int    `json:"asn"`
	AsOrganization string `json:"asOrganization"`
	Colo           string `json:"colo"`
	Country        string `json:"country"`
	City           string `json:"city"`
	Region         string `json:"region"`
	PostalCode     string `json:"postalCode"`
	Latitude       string `json:"latitude"`
	Longitude      string `json:"longitude"`
}
