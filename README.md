# github.com/bruceharrison1984/cloudflare-speed-test

This package mimics the behavior of the Cloudflare Speed Test webpage. It will run tests, then return values across channels.

## Basic Usage

```golang
import (
	"github.com/bruceharrison1984/cloudflare-speed-test/engines"
	"github.com/bruceharrison1984/cloudflare-speed-test/types"
)

speedTestSummaryChannel := make(chan *types.SpeedTestSummary)
cloudflareMetadataResults := make(chan *types.CloudflareMetadata)
exitChannel := make(chan struct{})
errorChannel := make(chan error)

engine := engines.NewTestEngine(speedTestSummaryChannel, cloudflareMetadataResults, exitChannel, errorChannel)

logger, err := s.Logger(errorChannel)
if err != nil {
    panic(err)
}

go engine.RunSpeedTest(cf.ctx)

// your preferred logger
logger.Info("Starting Speed Test")

var measurementId string
for {
    select {
    case _, ok := <-cloudflareMetadataResults:
        {
            if ok {
                // payload := types.NewUptimeMeasurementRequest(*metadata)
                // measurementResponse := commsUpClient.PostUptimeMeasurement(*payload)
                // measurementId = measurementResponse.Data.Id
                logger.Infof("Measurement Id: %s", measurementId)
                // if cf.verbose {
                // 	jsonObj, _ := json.MarshalIndent(metadata, "", "  ")
                // 	logger.Infof(string(jsonObj))
                // }
            }
        }
    case summary, ok := <-speedTestSummaryChannel:
        {
            if ok {
                if cf.verbose {
                    // jsonObj, _ := json.MarshalIndent(summary, "", "  ")
                    logger.Infof("Download: %+v\n", summary.Bandwidth.DownloadSpeedMbps)
                    logger.Infof("  Upload: %+v\n", summary.Bandwidth.UploadSpeedMbps)
                    logger.Infof("    Ping: %+v\n", summary.Bandwidth.Ping)
                }
            }
        }
    case _, ok := <-exitChannel:
        {
            if !ok {
                logger.Info("Speed test completed")
                return
            }
        }
    case err, ok := <-errorChannel:
        {
            if ok {
                logger.Error("Error occured, speed test failed!")
                logger.Error(err)
                return
            }
        }
    }
}
```
