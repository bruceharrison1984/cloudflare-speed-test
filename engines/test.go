package engines

import (
	"context"
	"math/rand"
	"net/http"
	"time"

	"github.com/bruceharrison1984/cloudflare-speed-test/aggregators"
	"github.com/bruceharrison1984/cloudflare-speed-test/clients"
	"github.com/bruceharrison1984/cloudflare-speed-test/providers"
	"github.com/bruceharrison1984/cloudflare-speed-test/types"
)

type ISpeedTestEngine interface {
	RunSpeedTest(ctx context.Context, testCases []types.SpeedTestCase)
}

/*
The primary engine for running CloudFlare speed tests.

This is probably the one you want.
*/
type cloudflareSpeedTestEngine struct {
	SpeedTestSummaryChannel   chan *types.SpeedTestSummary   // Piping this channel will give access to the final summary once a run completes
	CloudflareMetadataResults chan *types.CloudflareMetadata // Listen here for test metadata
	ExitChannel               chan struct{}                  // Listen here to end the listener loop
	ErrorChannel              chan error                     // Errors are reported here, which also ends the listener loop
}

/* Create a new test engine */
func NewTestEngine(
	speedTestSummaryChannel chan *types.SpeedTestSummary,
	exitChannel chan struct{},
	errorChannel chan error) ISpeedTestEngine {
	return &cloudflareSpeedTestEngine{
		SpeedTestSummaryChannel: speedTestSummaryChannel, // this should be passed in
		ExitChannel:             exitChannel,             // this should be passed in
		ErrorChannel:            errorChannel,            // this should be passed in
	}
}

/* Run the speed tests */
func (t *cloudflareSpeedTestEngine) RunSpeedTest(ctx context.Context, testCases []types.SpeedTestCase) {

	testId := rand.Int63()
	httpTestClient := &http.Client{
		Timeout:   time.Second * 20,
		Transport: clients.NewCloudflareSpeedTestTransport(),
	}

	// testConfig, iterations := config.GetDefaultConfig()
	var iterations int
	for i := 0; i < len(testCases); i++ {
		iterations += testCases[i].Iterations
	}

	rawBandwidthResultsChan := make(chan *types.RawBandwidthClientResult, iterations)

	urlProvider := providers.NewUrlProvider()
	metadataClient := clients.NewMetadataClient(httpTestClient, urlProvider)
	bandwidthEngine := NewBandwidthEngine(clients.NewBandwidthClient(httpTestClient), urlProvider)

	metadata, err := metadataClient.FetchMetadata()
	if err != nil {
		t.ErrorChannel <- err
		return
	}

	resultsEngine := aggregators.NewResultsAggregator(t.SpeedTestSummaryChannel, metadata)
	go resultsEngine.Listen(rawBandwidthResultsChan)
	bandwidthEngine.RunTest(ctx, testId, testConfig, rawBandwidthResultsChan, t.ErrorChannel)

	close(t.ExitChannel)
	// close(t.ErrorChannel) not sure if this should be closed?
	// close(t.SpeedTestSummaryChannel) not sure if this should be closed?
}
