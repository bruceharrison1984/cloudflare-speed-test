package engines

import (
	"math"
	"time"

	"github.com/bruceharrison1984/cloudflare-speed-test/types"
	"github.com/montanaflynn/stats"
)

/* Engine that aggregates the test results */
type ResultsEngine struct {
	speedTestMetadata       *types.CloudflareMetadata
	SpeedTestSummaryChannel chan *types.SpeedTestSummary   // Piping this channel will give access to the final summary once a run completes
	bandwidthResults        []*types.BandwidthClientResult // Internal array of results
	ErrorChannel            chan error                     // Piping this channel will give access to real-time errors
}

/* Create a new results engine */
func NewResultsEngine(summaryChannel chan *types.SpeedTestSummary, metadata *types.CloudflareMetadata, outputErrors chan error) *ResultsEngine {
	return &ResultsEngine{SpeedTestSummaryChannel: summaryChannel, speedTestMetadata: metadata, ErrorChannel: outputErrors}
}

/*
Listen to the raw data channels and compile metrics in real-time based on the results.

Compiled results are available on the SpeedTestSummaryChannel.
*/
func (engine *ResultsEngine) Listen(rawResultsChan chan *types.RawBandwidthClientResult, errorChan chan error) {
	for {
		select {
		case rawResult, ok := <-rawResultsChan:
			{
				if !ok {
					return
				}

				calculatedResult := engine.calculateMetrics(rawResult.ResultType, rawResult.ServerTiming, rawResult.TTFB, rawResult.TTLB, rawResult.PayloadSizeBytes)
				engine.bandwidthResults = append(engine.bandwidthResults, &calculatedResult)
				engine.SpeedTestSummaryChannel <- &types.SpeedTestSummary{
					TestResults: engine.bandwidthResults,
					Metadata:    engine.speedTestMetadata,
					Bandwidth:   engine.CalculatePercentiles(),
				}
			}
		case err, ok := <-errorChan:
			{
				if ok {
					engine.ErrorChannel <- err
					return
				}
			}
		}
	}
}

func (engine ResultsEngine) calculateMetrics(testType types.BandwidthTestType, serverTiming time.Duration, ttfb time.Duration, ttlb time.Duration, responseSizeBytes int64) types.BandwidthClientResult {
	ping := (ttfb - serverTiming).Seconds()
	if ping <= 0 {
		ping = (time.Millisecond * 1).Seconds()
	}

	transferDuration := ping + ttlb.Seconds()
	responseSizeMegabits := float64(responseSizeBytes) / 125000

	speedMbps := responseSizeMegabits / transferDuration

	return types.BandwidthClientResult{
		ResultType:       testType,
		PayloadSizeBytes: responseSizeBytes,
		ServerTiming:     Round(serverTiming.Seconds()),
		Ping:             Round(ping),
		TransferDuration: Round(transferDuration),
		SpeedMbps:        Round(speedMbps),
		Ttfb:             Round(ttfb.Seconds()),
		Ttlb:             Round(ttlb.Seconds()),
	}
}

func Round(val float64) float64 {
	return math.Round(val*(math.Pow10(3))) / math.Pow10(3)
}

func (engine ResultsEngine) CalculatePercentiles() *types.BandwidthClientResultSummary {
	// summary percentile
	var pings []float64
	var downloads []float64
	var uploads []float64

	for i := 0; i < len(engine.bandwidthResults); i++ {
		pings = append(pings, engine.bandwidthResults[i].Ping)
		if engine.bandwidthResults[i].ResultType == types.Download {
			downloads = append(downloads, engine.bandwidthResults[i].SpeedMbps)
		} else {
			uploads = append(uploads, engine.bandwidthResults[i].SpeedMbps)
		}
	}

	percentilePing, _ := stats.PercentileNearestRank(pings, 90)
	percentileDownload, _ := stats.PercentileNearestRank(downloads, 90)
	percentileUpload, _ := stats.PercentileNearestRank(uploads, 90)

	return &types.BandwidthClientResultSummary{
		Ping:              percentilePing,
		DownloadSpeedMbps: percentileDownload,
		UploadSpeedMbps:   percentileUpload,
	}
}
