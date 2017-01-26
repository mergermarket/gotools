package tools

import (
	"net/http"
	"github.com/felixge/httpsnoop"
	"strconv"
	"fmt"
)

type logger interface {
	Debug(...interface{})
	Info(...interface{})
	Error(...interface{})
}

// WrapWithTelemetry takes your http.Handler and adds debug logs with request details and marks metrics
func WrapWithTelemetry(routeName string, router http.Handler, logger logger, statsd StatsD) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Debug(r.Method, "at", r.URL.String())
		metrics := httpsnoop.CaptureMetrics(router, w, r)

		logResult(routeName, metrics, statsd, logger, r.URL.String())
	})
}

func logResult(routeName string, metrics httpsnoop.Metrics, statsd StatsD, logger logger, url string) {
	statsd.Histogram("web.response_time", float64(metrics.Duration.Nanoseconds())/1000000, "route:"+routeName, "response:" + strconv.Itoa(metrics.Code))

	message := fmt.Sprint("Request to ", url, " had response code ", metrics.Code)
	if metrics.Code >= 400 {
		logger.Error(message)
	} else {
		logger.Info(message)
	}
}
