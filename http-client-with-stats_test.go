package tools

import (
	"fmt"
	"github.com/mergermarket/gotools/testtools"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type fakeClock struct {
	times chan time.Time
}

func (f *fakeClock) Now() time.Time {
	return <-f.times
}

func TestTelemetryHTTPClient_Do(t *testing.T) {

	fc := &fakeClock{make(chan time.Time, 2)}
	callee := "my-remote-service"
	msd := &testtools.MockStatsD{}
	hc := http.DefaultClient
	operationMapper := func(r *http.Request) (operationTag string) {
		return "my-operation"
	}
	wc := &httpClientWithStats{statsd: msd, httpClient: hc, clock: fc, callee: callee, operationTagFunc: operationMapper}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello World")
	}))
	defer ts.Close()

	start := time.Now()
	fc.times <- start
	fc.times <- start.Add(100 * time.Millisecond)

	req, _ := http.NewRequest("GET", ts.URL, nil)

	assert.Len(t, msd.Calls, 0)

	resp, _ := wc.Do(req)

	expectedTags := []string{"http_callee:my-remote-service", "method:GET", "operation:my-operation", "resp_status:200"}

	assert.Len(t, msd.Calls, 2)
	assert.NotNil(t, resp)
	assert.Equal(t, "Histogram", msd.Calls[0].Method)
	assert.Equal(t, 100.0, msd.Calls[0].Args.Value)
	assert.Equal(t, expectedTags, msd.Calls[0].Args.Tags)
	assert.Equal(t, "Incr", msd.Calls[1].Method)
	assert.Equal(t, "http_client.response_success", msd.Calls[1].Args.Name)
	assert.Equal(t, expectedTags, msd.Calls[1].Args.Tags)
}

func TestTelemetryHTTPClient_Do_Error(t *testing.T) {

	fc := &fakeClock{make(chan time.Time, 2)}

	msd := &testtools.MockStatsD{}
	hc := http.DefaultClient
	wc := &httpClientWithStats{statsd: msd, httpClient: hc, clock: fc, callee: "my-remote-service"}

	start := time.Now()
	fc.times <- start
	fc.times <- start.Add(100 * time.Millisecond)

	req, _ := http.NewRequest("GET", "http://not-a-domain.zyzyuyziuy", nil)

	assert.Len(t, msd.Calls, 0)

	resp, err := wc.Do(req)

	assert.Len(t, msd.Calls, 1)
	assert.NotNil(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, "Incr", msd.Calls[0].Method)
	assert.Equal(t, "http_client.response_error", msd.Calls[0].Args.Name)
	assert.Equal(t, []string{"http_callee:my-remote-service", "method:GET"}, msd.Calls[0].Args.Tags)
}

func TestTelemetryHTTPClient_Get(t *testing.T) {

	fc := &fakeClock{make(chan time.Time, 2)}

	msd := &testtools.MockStatsD{}
	hc := http.DefaultClient
	wc := &httpClientWithStats{statsd: msd, httpClient: hc, clock: fc, callee: "my-remote-service"}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello World")
	}))
	defer ts.Close()

	start := time.Now()
	fc.times <- start
	fc.times <- start.Add(100 * time.Millisecond)

	assert.Len(t, msd.Calls, 0)

	resp, _ := wc.Get(ts.URL)

	assert.Len(t, msd.Calls, 2)
	assert.NotNil(t, resp)
	assert.Equal(t, "Histogram", msd.Calls[0].Method)
	assert.Equal(t, 100.0, msd.Calls[0].Args.Value)
	assert.Equal(t, []string{"http_callee:my-remote-service", "method:GET", "resp_status:200"}, msd.Calls[0].Args.Tags)

}

func TestTelemetryHTTPClient_Post(t *testing.T) {
	fc := &fakeClock{make(chan time.Time, 2)}

	msd := &testtools.MockStatsD{}
	hc := http.DefaultClient
	wc := &httpClientWithStats{statsd: msd, httpClient: hc, clock: fc, callee: "my-remote-service"}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello World")
	}))
	defer ts.Close()

	start := time.Now()
	fc.times <- start
	fc.times <- start.Add(100 * time.Millisecond)

	assert.Len(t, msd.Calls, 0)

	resp, _ := wc.Post(ts.URL, "application/json", strings.NewReader(`{"hello":"world"}`))

	assert.Len(t, msd.Calls, 2)
	assert.NotNil(t, resp)
	assert.Equal(t, "Histogram", msd.Calls[0].Method)
	assert.Equal(t, 100.0, msd.Calls[0].Args.Value)
	assert.Equal(t, []string{"http_callee:my-remote-service", "method:POST", "resp_status:200"}, msd.Calls[0].Args.Tags)

}
