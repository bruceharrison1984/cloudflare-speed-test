package clients

import (
	"bytes"
	"context"
	"io"
	"math"
	"net/http"
	"net/http/httptrace"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bruceharrison1984/cloudflare-speed-test/types"
)

/* This is the client that is used during the bandwidth test */
type BandwidthTestClient struct {
	Http *http.Client
}

/* Create a new bandwidth client */
func NewBandwidthClient(http *http.Client) *BandwidthTestClient {
	return &BandwidthTestClient{http}
}

/* Begin running the bandwidth test */
func (client BandwidthTestClient) RunTest(ctx context.Context, url string, testId int64, payloadLength int64, testType types.BandwidthTestType) (*types.RawBandwidthClientResult, error) {
	var handshakeComplete time.Time
	var ttfb, ttlb time.Duration

	trace := &httptrace.ClientTrace{
		GotConn: func(gci httptrace.GotConnInfo) {
			// retreiving a connection can be slow, so we don't start timing until it has completed
			handshakeComplete = time.Now()
		},
		GotFirstResponseByte: func() {
			ttfb = time.Since(handshakeComplete)
		},
		PutIdleConn: func(err error) {
			ttlb = time.Since(handshakeComplete)
		},
	}

	requestBody := client.createRequestBody(testType, payloadLength)

	requestType := http.MethodGet
	if testType == types.Upload {
		requestType = http.MethodPost
	}

	req, err := http.NewRequestWithContext(httptrace.WithClientTrace(ctx, trace), requestType, url, requestBody)
	if err != nil {
		return nil, err
	}

	resp, err := client.Http.Transport.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	loggingToken := string(body)
	if testType == types.Download {
		loggingToken = strings.Split(string(body), "___")[1] /// logging token is found at the end of a download request
	}

	payloadSize := int64(len(body))
	if testType == types.Upload {
		payloadSize = payloadLength
	}

	rawResults := &types.RawBandwidthClientResult{ /// switch output based on test type
		ServerTiming:     client.getServerTiming(&resp.Header),
		ResultType:       testType,
		TTFB:             ttfb,
		TTLB:             ttlb,
		PayloadSizeBytes: payloadSize, // depends on test type
		LoggingToken:     loggingToken,
	}

	return rawResults, nil
}

/* Create a zero'd out payload for use in bandwidth test */
func (runner BandwidthTestClient) createRequestBody(testType types.BandwidthTestType, payloadLength int64) io.Reader {
	if testType == types.Download {
		return nil
	}
	rawBody := bytes.Repeat([]byte{0x30}, int(payloadLength))
	return strings.NewReader(string(rawBody))
}

// Extract server time from response headers
func (runner BandwidthTestClient) getServerTiming(headers *http.Header) time.Duration {
	var rawTiming = headers.Get("server-timing")

	regex, _ := regexp.Compile("dur=([0-9.]+)")
	var matches = regex.FindAllStringSubmatch(rawTiming, 2)

	i, _ := strconv.ParseFloat(matches[0][1], 32)
	return time.Duration(math.Round(i)) * time.Millisecond
}
