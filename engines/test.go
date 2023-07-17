package engines

import (
	"context"
	"math/rand"
	"net/http"
	"time"

	"github.com/bruceharrison1984/cloudflare-speed-test/aggregators"
	"github.com/bruceharrison1984/cloudflare-speed-test/clients"
	"github.com/bruceharrison1984/cloudflare-speed-test/config"
	"github.com/bruceharrison1984/cloudflare-speed-test/providers"
	"github.com/bruceharrison1984/cloudflare-speed-test/types"
)

type ISpeedTestEngine interface {
	RunSpeedTest(ctx context.Context)
}

/*
The primary engine for running CloudFlare speed tests.

This is probably the one you want.
*/
type cloudflareSpeedTestEngine struct {
	SpeedTestSummaryChannel   chan *types.SpeedTestSummary   // Piping this channel will give access to the final summary once a run completes
	CloudflareMetadataResults chan *types.CloudflareMetadata // Listen here for test metadata
	Exit                      chan struct{}                  // Listen here to end the listener loop
	Errors                    chan error                     // Errors are reported here, which also ends the listener loop
}

/* Create a new test engine */
func NewTestEngine(
	speedTestSummaryChannel chan *types.SpeedTestSummary,
	cloudflareMetadataResults chan *types.CloudflareMetadata,
	exitChannel chan struct{},
	errorChannel chan error) ISpeedTestEngine {
	return &cloudflareSpeedTestEngine{
		SpeedTestSummaryChannel:   speedTestSummaryChannel,   // this should be passed in
		CloudflareMetadataResults: cloudflareMetadataResults, // this should be passed in
		Exit:                      exitChannel,               // this should be passed in
		Errors:                    errorChannel,              // this should be passed in
	}
}

/* Run the speed tests */
func (t *cloudflareSpeedTestEngine) RunSpeedTest(ctx context.Context) {

	testId := rand.Int63()
	httpTestClient := &http.Client{
		Timeout:   time.Second * 20,
		Transport: clients.NewCloudflareSpeedTestTransport(),
	}

	testConfig, iterations := config.GetDefaultConfig()

	rawBandwidthResultsChan := make(chan *types.RawBandwidthClientResult, iterations)

	urlProvider := providers.NewUrlProvider()
	metadataClient := clients.NewMetadataClient(httpTestClient, urlProvider)
	bandwidthEngine := NewBandwidthEngine(clients.NewBandwidthClient(httpTestClient), urlProvider)

	metadata, err := metadataClient.FetchMetadata()
	if err != nil {
		t.Errors <- err
		return
	}

	resultsEngine := aggregators.NewResultsAggregator(t.SpeedTestSummaryChannel, metadata, t.Errors)
	go resultsEngine.Listen(rawBandwidthResultsChan, t.Errors)
	bandwidthEngine.RunTest(ctx, testId, testConfig, rawBandwidthResultsChan, t.Errors)

	close(t.Exit)
	close(t.Errors)
	close(t.SpeedTestSummaryChannel)
}
