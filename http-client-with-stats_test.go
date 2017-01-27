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

func TestHTTPClientWithStats_Do(t *testing.T) {

	fc := &fakeClock{make(chan time.Time, 2)}
	msd := &testtools.MockStatsD{}
	hc := http.DefaultClient
	wc := &httpClientWithStats{statsd: msd, httpClient: hc, clock: fc}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello World")
	}))
	defer ts.Close()

	start := time.Now()
	fc.times <- start
	fc.times <- start.Add(100 * time.Millisecond)

	req, _ := http.NewRequest("GET", ts.URL, nil)

	assert.Len(t, msd.Calls, 0)

	resp, _ := wc.Do(req, "http_callee:my-remote-service", "operation:my-operation")

	expectedTags := []string{"http_callee:my-remote-service", "operation:my-operation", "method:GET", "resp_status:200"}

	assert.Len(t, msd.Calls, 2)
	assert.NotNil(t, resp)
	assert.Equal(t, "Histogram", msd.Calls[0].Method)
	assert.Equal(t, 100.0, msd.Calls[0].Args.Value)
	assert.Equal(t, expectedTags, msd.Calls[0].Args.Tags)
	assert.Equal(t, "Incr", msd.Calls[1].Method)
	assert.Equal(t, "http_client.response_success", msd.Calls[1].Args.Name)
	assert.Equal(t, expectedTags, msd.Calls[1].Args.Tags)
}

func TestHTTPClientWithStats_Do_Error(t *testing.T) {

	fc := &fakeClock{make(chan time.Time, 2)}

	msd := &testtools.MockStatsD{}
	hc := http.DefaultClient
	wc := &httpClientWithStats{statsd: msd, httpClient: hc, clock: fc}

	start := time.Now()
	fc.times <- start
	fc.times <- start.Add(100 * time.Millisecond)

	req, _ := http.NewRequest("GET", "http://not-a-domain.zyzyuyziuy", nil)

	assert.Len(t, msd.Calls, 0)

	resp, err := wc.Do(req, "http_callee:my-remote-service", "operation:my-operation")

	assert.Len(t, msd.Calls, 1)
	assert.NotNil(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, "Incr", msd.Calls[0].Method)
	assert.Equal(t, "http_client.response_error", msd.Calls[0].Args.Name)
	assert.Equal(t, []string{"http_callee:my-remote-service", "operation:my-operation", "method:GET"}, msd.Calls[0].Args.Tags)
}

func TestHTTPClientWithStats_Get(t *testing.T) {

	fc := &fakeClock{make(chan time.Time, 2)}

	msd := &testtools.MockStatsD{}
	hc := http.DefaultClient
	wc := &httpClientWithStats{statsd: msd, httpClient: hc, clock: fc}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello World")
	}))
	defer ts.Close()

	start := time.Now()
	fc.times <- start
	fc.times <- start.Add(100 * time.Millisecond)

	assert.Len(t, msd.Calls, 0)

	resp, _ := wc.Get(ts.URL, "http_callee:my-remote-service", "operation:my-operation")

	assert.Len(t, msd.Calls, 2)
	assert.NotNil(t, resp)
	assert.Equal(t, "Histogram", msd.Calls[0].Method)
	assert.Equal(t, 100.0, msd.Calls[0].Args.Value)
	assert.Equal(t, []string{"http_callee:my-remote-service", "operation:my-operation", "method:GET", "resp_status:200"}, msd.Calls[0].Args.Tags)

}

func TestHttpClientWithStats_Post(t *testing.T) {
	fc := &fakeClock{make(chan time.Time, 2)}

	msd := &testtools.MockStatsD{}
	hc := http.DefaultClient
	wc := &httpClientWithStats{statsd: msd, httpClient: hc, clock: fc}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello World")
	}))
	defer ts.Close()

	start := time.Now()
	fc.times <- start
	fc.times <- start.Add(100 * time.Millisecond)

	assert.Len(t, msd.Calls, 0)

	resp, _ := wc.Post(ts.URL, "application/json", strings.NewReader(`{"hello":"world"}`), "http_callee:my-remote-service", "operation:my-operation")

	assert.Len(t, msd.Calls, 2)
	assert.NotNil(t, resp)
	assert.Equal(t, "Histogram", msd.Calls[0].Method)
	assert.Equal(t, 100.0, msd.Calls[0].Args.Value)
	assert.Equal(t, []string{"http_callee:my-remote-service", "operation:my-operation", "method:POST", "resp_status:200"}, msd.Calls[0].Args.Tags)

}
