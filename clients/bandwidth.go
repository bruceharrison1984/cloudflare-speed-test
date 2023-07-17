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

type BandwidthClient struct {
	Http *http.Client
}

func NewBandwidthClient(http *http.Client) *BandwidthClient {
	return &BandwidthClient{http}
}

func (client BandwidthClient) RunTest(ctx context.Context, url string, testId int64, payloadLength int64, testType types.BandwidthTestType) (*types.RawBandwidthClientResult, error) {
	var handshakeComplete time.Time
	var ttfb, ttlb time.Duration

	trace := &httptrace.ClientTrace{
		GotConn: func(gci httptrace.GotConnInfo) {
			handshakeComplete = time.Now()
		}, // retreiving a connection can be slow, so we don't start timing until it has completed
		GotFirstResponseByte: func() {
			ttfb = time.Since(handshakeComplete)
		},
		PutIdleConn: func(err error) {
			ttlb = time.Since(handshakeComplete)
		},
	}

	requestBody := client.getRequestBody(testType, payloadLength)

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

func (runner BandwidthClient) getRequestBody(testType types.BandwidthTestType, payloadLength int64) io.Reader {
	if testType == types.Download {
		return nil
	}
	rawBody := bytes.Repeat([]byte{0x30}, int(payloadLength))
	return strings.NewReader(string(rawBody))
}

// Extract server time from response headers
func (runner BandwidthClient) getServerTiming(headers *http.Header) time.Duration {
	var rawTiming = headers.Get("server-timing")

	regex, _ := regexp.Compile("dur=([0-9.]+)")
	var matches = regex.FindAllStringSubmatch(rawTiming, 2)

	i, _ := strconv.ParseFloat(matches[0][1], 32)
	return time.Duration(math.Round(i)) * time.Millisecond
}
