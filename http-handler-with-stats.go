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

		logResult(routeName, metrics, statsd, logger, r)
	})
}

func logResult(routeName string, metrics httpsnoop.Metrics, statsd StatsD, logger logger, req *http.Request) {
	responseTag := fmt.Sprintf("response:%d", metrics.Code)
	statsd.Histogram(WebResponseTimeKey, float64(metrics.Duration.Nanoseconds())/1000000, "route:"+routeName, responseTag)
	responseCodeKey := fmt.Sprintf(WebResponseCodeFormatKey, metrics.Code)
	tags := []string{"route:" + routeName, responseTag}
	caller := req.Header.Get("X-Component")
	if caller != "" {
		tags = append(tags, "caller:"+caller)
	}
	statsd.Incr(responseCodeKey, tags...)
	statsd.Incr(WebResponseCodeAllKey, tags...)

	message := fmt.Sprint("Request to ", req.URL.String(), " had response code ", metrics.Code)
	if metrics.Code >= 400 {
		logger.Error(message)
	} else {
		logger.Info(message) // todo should this be debug?
	}
}
