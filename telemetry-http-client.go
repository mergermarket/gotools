package tools

import (
	"net/http"
	"time"
	"fmt"
	"io"
)

type TelemetryHttpClient interface {
	Do(r *http.Request) (*http.Response, error)
	Get(url string) (*http.Response, error)
	Post(url string, bodyType string, body io.Reader) (*http.Response, error)
}

type clock interface {
	Now() time.Time
}

type telemetryHTTPClient struct {
	httpClient *http.Client
	statsd StatsD
	clock clock
	callee string
}

const responseTimeKey = "http_client.response_time_ms"
const responseErrorKey = "http_client.response_error"
const responseSuccessKey = "http_client.response_success"

func (thc *telemetryHTTPClient) Do(r *http.Request) (*http.Response, error) {
	start := thc.clock.Now()
	resp, err := thc.httpClient.Do(r)
	if err != nil {
		thc.statsd.Incr(responseErrorKey, "http_callee:"+thc.callee, "method:"+r.Method)
	} else {
		finish := thc.clock.Now()
		duration := (finish.Nanosecond()-start.Nanosecond())/1000000
		thc.statsd.Histogram(responseTimeKey, float64(duration), "http_callee:"+thc.callee, "method:"+r.Method, fmt.Sprintf("resp_status:%d",resp.StatusCode))
		thc.statsd.Incr(responseSuccessKey, "http_callee:"+thc.callee, "method:"+r.Method)
	}
	return resp, err
}

func (thc *telemetryHTTPClient) Get(url string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return resp, err
	}
	return thc.Do(req)
}

func (thc *telemetryHTTPClient) Post(url string, bodyType string, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return resp, err
	}
	req.Header.Set("Content-Type", bodyType)
	return thc.Do(req)
}

type timeClock struct {}

func (c *timeClock) Now() time.Time {
	return time.Now()
}

func TelemetryHTTPClient(client *http.Client, statsd StatsD, callee string) TelemetryHttpClient {
	return &telemetryHTTPClient{statsd:statsd, httpClient:client, clock:&timeClock{}, callee:callee}
}
