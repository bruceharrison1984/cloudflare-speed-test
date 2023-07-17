package engines

import (
	"context"

	"github.com/bruceharrison1984/cloudflare-speed-test/clients"
	"github.com/bruceharrison1984/cloudflare-speed-test/providers"
	"github.com/bruceharrison1984/cloudflare-speed-test/types"
)

/* This engine runs the bandwidth test */
type BandwidthEngine struct {
	SpeedTestClient *clients.BandwidthClient
	urlProvider     providers.UrlProvider
}

/* Create a new bandwidth engine */
func NewBandwidthEngine(speedTestClient *clients.BandwidthClient, urlProvider providers.UrlProvider) *BandwidthEngine {
	return &BandwidthEngine{speedTestClient, urlProvider}
}

/* Run the bandwidth test */
func (engine *BandwidthEngine) RunTest(ctx context.Context, testId int64, testCases []types.SpeedTestCase, rawResultsChannel chan *types.RawBandwidthClientResult, errorChan chan error) {
	for x := 0; x < len(testCases); x++ {
		testCase := testCases[x]
		for i := 0; i < testCase.Iterations; i++ {

			url := engine.urlProvider.GetDownloadTestUrl(testId, testCase.PayloadSize)
			if testCase.TestType == types.Upload {
				url = engine.urlProvider.GetUploadUrl(testId)
			}

			result, err := engine.SpeedTestClient.RunTest(ctx, url, testId, testCase.PayloadSize, testCase.TestType)
			if err != nil {
				errorChan <- err
				close(rawResultsChannel)
				close(errorChan)
				return
			}
			rawResultsChannel <- result
		}
	}
	close(rawResultsChannel)
}
