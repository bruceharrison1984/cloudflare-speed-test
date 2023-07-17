# github.com/bruceharrison1984/cloudflare-speed-test

This package mimics the behavior of the Cloudflare Speed Test webpage. It will run tests, then return values across channels.

Each subsequent test iteration will output the aggregated results on the `types.SpeedTestSummary` channel.

The `types.SpeedTestSummary` object is as follows:
| Property | Description |
|---|---|
| Metadata | Information about the client connection |
| Bandwidth | Current aggregate test results (90th percentile)|
| TestResults | An array of the raw test results. New results are appended to this array for each test iteration |

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

for {
    select {
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

Running this should output something similar to:

```sh
...
I: 13:33:18     Ping: 0.014
I: 13:33:19 Download: 673.508
I: 13:33:19   Upload: 385.959
I: 13:33:19     Ping: 0.014
I: 13:33:19 Download: 673.508
I: 13:33:19   Upload: 413.899
I: 13:33:19     Ping: 0.014
I: 13:33:20 Download: 673.508
I: 13:33:20   Upload: 413.899
I: 13:33:20     Ping: 0.013
I: 13:33:21 Download: 673.508
I: 13:33:21   Upload: 413.899
I: 13:33:21     Ping: 0.013
I: 13:33:22 Speed test completed
```
