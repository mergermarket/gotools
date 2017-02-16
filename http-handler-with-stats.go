package tools

import (
	"fmt"
	"github.com/felixge/httpsnoop"
	"net/http"
)

type logger interface {
	Debug(...interface{})
	Info(...interface{})
	Error(...interface{})
}

// HTTPHandlerWithStats takes an http.Handler and adds the sending of response time metrics to DataDog, and debug logging of request details
func HTTPHandlerWithStats(routeName string, router http.Handler, logger logger, statsd StatsD) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Debug(r.Method, "at", r.URL.String())
		metrics := httpsnoop.CaptureMetrics(router, w, r)

		logResult(routeName, metrics, statsd, logger, r.URL.String())
	})
}

func logResult(routeName string, metrics httpsnoop.Metrics, statsd StatsD, logger logger, url string) {
	responseTag := fmt.Sprintf("response:%d", metrics.Code)
	statsd.Histogram(WebResponseTimeKey, float64(metrics.Duration.Nanoseconds())/1000000, "route:"+routeName, responseTag)
	responseCodeKey := fmt.Sprintf(WebResponseCodeFormatKey, metrics.Code)
	statsd.Incr(responseCodeKey, "route:"+routeName, responseTag)
	statsd.Incr(WebResponseCodeAllKey, "route:"+routeName, responseTag)

	message := fmt.Sprint("Request to ", url, " had response code ", metrics.Code)
	if metrics.Code >= 400 {
		logger.Error(message)
	} else {
		logger.Info(message)
	}
}
