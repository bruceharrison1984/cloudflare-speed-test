package engines

import (
	"context"
	"math/rand"
	"net/http"
	"time"

	"github.com/bruceharrison1984/cloudflare-speed-test/clients"
	"github.com/bruceharrison1984/cloudflare-speed-test/config"
	"github.com/bruceharrison1984/cloudflare-speed-test/providers"
	"github.com/bruceharrison1984/cloudflare-speed-test/types"
)

/*
The primary engine for running CloudFlare speed tests.

This is probably the one you want.
*/
type CloudflareSpeedTestEngine struct {
	SpeedTestSummaryChannel   chan *types.SpeedTestSummary   // Piping this channel will give access to the final summary once a run completes
	CloudflareMetadataResults chan *types.CloudflareMetadata // Listen here for test metadata
	Exit                      chan struct{}                  // Listen here to end the listener loop
	Errors                    chan error                     // Errors are reported here, which also ends the listener loop
}

/* Create a new test engine */
func NewTestEngine() *CloudflareSpeedTestEngine {
	return &CloudflareSpeedTestEngine{
		SpeedTestSummaryChannel:   make(chan *types.SpeedTestSummary),   // this should be passed in
		CloudflareMetadataResults: make(chan *types.CloudflareMetadata), // this should be passed in
		Exit:                      make(chan struct{}),                  // this should be passed in
		Errors:                    make(chan error),                     // this should be passed in
	}
}

/* Run the speed tests */
func (t *CloudflareSpeedTestEngine) RunSpeedTest(ctx context.Context) {

	testId := rand.Int63()
	httpTestClient := &http.Client{
		Timeout:   time.Second * 20,
		Transport: &clients.CloudflareSpeedTestTransport{}}

	testConfig, iterations := config.GetDefaultConfig()

	rawBandwidthResultsChan := make(chan *types.RawBandwidthClientResult, iterations)

	urlProvider := providers.UrlProvider{}
	metadataClient := clients.NewMetadataClient(httpTestClient, urlProvider)
	bandwidthEngine := NewBandwidthEngine(clients.NewBandwidthClient(httpTestClient), urlProvider)

	metadata, err := metadataClient.FetchMetadata()
	if err != nil {
		t.Errors <- err
		return
	}

	resultsEngine := NewResultsEngine(t.SpeedTestSummaryChannel, metadata, t.Errors)
	go resultsEngine.Listen(rawBandwidthResultsChan, t.Errors)
	bandwidthEngine.RunTest(ctx, testId, testConfig, rawBandwidthResultsChan, t.Errors)

	close(t.Exit)
	close(t.Errors)
	close(t.SpeedTestSummaryChannel)
}
