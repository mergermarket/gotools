package tools

import (
	"github.com/mergermarket/gotools/testtools"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func checkTimingMetricCalled(t *testing.T, statsd *testtools.MockStatsD, routeName string, response int) {
	call, err := statsd.Call()

	assert.Nil(t, err, "No call made to MockStatsD")
	assert.Equal(t, "Histogram", call.Method)

	assert.Equal(t, "web.response_time", call.Args.Name, "Expected name of metric to be 'web.response_time', but got:", call.Args.Name)

	assert.Contains(t, call.Args.Tags, "route:"+routeName)
	assert.Contains(t, call.Args.Tags, "response:"+strconv.Itoa(response))
}

type MockHandler struct {
	response int
}

func (h MockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(h.response)
}

func TestHTTPHandlerWithStats(t *testing.T) {
	statsd := &testtools.MockStatsD{}
	logger := &testtools.MockLogger{}
	router := &MockHandler{response: http.StatusOK}
	httpHandler := HTTPHandlerWithStats("route", router, logger, statsd)

	httpHandler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "http://example.com", nil))

	lastLogCall := logger.LastCall()

	assert.NotNil(t, lastLogCall, "Expected call to be made")
	assert.Equal(t, "Info", lastLogCall.Method)
	assert.Equal(t, "[Request to http://example.com had response code 200]", lastLogCall.Args.Msg)

	checkTimingMetricCalled(t, statsd, "route", http.StatusOK)
}

func TestHTTPHandlerWithStats_Error(t *testing.T) {
	statsd := &testtools.MockStatsD{}
	logger := &testtools.MockLogger{}
	router := &MockHandler{response: http.StatusInternalServerError}
	httpHandler := HTTPHandlerWithStats("route", router, logger, statsd)

	httpHandler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "http://example.com", nil))

	lastLogCall := logger.LastCall()

	assert.NotNil(t, lastLogCall, "Expected call to be made")
	assert.Equal(t, "Error", lastLogCall.Method)
	assert.Equal(t, "[Request to http://example.com had response code 500]", lastLogCall.Args.Msg)

	checkTimingMetricCalled(t, statsd, "route", http.StatusInternalServerError)
}
