package tools

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type HTTPClientWithStats interface {
	Do(r *http.Request) (*http.Response, error)
	Get(url string) (*http.Response, error)
	Post(url string, bodyType string, body io.Reader) (*http.Response, error)
}

type clock interface {
	Now() time.Time
}

type httpClientWithStats struct {
	httpClient       *http.Client
	statsd           StatsD
	clock            clock
	callee           string
	operationTagFunc operationTagFunction
}

type operationTagFunction func(r *http.Request) (operationTag string)

const responseTimeKey = "http_client.response_time_ms"
const responseErrorKey = "http_client.response_error"
const responseSuccessKey = "http_client.response_success"

func (thc *httpClientWithStats) Do(r *http.Request) (*http.Response, error) {
	tags := []string{"http_callee:" + thc.callee, "method:" + r.Method}
	if thc.operationTagFunc != nil {
		tags = append(tags, fmt.Sprintf("operation:%s", thc.operationTagFunc(r)))
	}
	start := thc.clock.Now()
	resp, err := thc.httpClient.Do(r)
	if err != nil {
		thc.statsd.Incr(responseErrorKey, tags...)
	} else {
		finish := thc.clock.Now()
		duration := (finish.Nanosecond() - start.Nanosecond()) / 1000000
		tags = append(tags, fmt.Sprintf("resp_status:%d", resp.StatusCode))
		thc.statsd.Histogram(responseTimeKey, float64(duration), tags...)
		thc.statsd.Incr(responseSuccessKey, tags...)
	}
	return resp, err
}

func (thc *httpClientWithStats) Get(url string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return resp, err
	}
	return thc.Do(req)
}

func (thc *httpClientWithStats) Post(url string, bodyType string, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return resp, err
	}
	req.Header.Set("Content-Type", bodyType)
	return thc.Do(req)
}

type timeClock struct{}

func (c *timeClock) Now() time.Time {
	return time.Now()
}

func NewHTTPClientWithStats(client *http.Client, statsd StatsD, callee string, operationTagDeterminer func(*http.Request) string) HTTPClientWithStats {
	return &httpClientWithStats{statsd: statsd, httpClient: client, clock: &timeClock{}, callee: callee, operationTagFunc: operationTagDeterminer}
}
