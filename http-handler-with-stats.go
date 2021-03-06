package tools

import (
	"fmt"
	"github.com/felixge/httpsnoop"
	"net/http"
)

// HTTPHandlerWithStats takes an http.Handler and adds the sending of response time metrics to DataDog, and debug logging of request details
func HTTPHandlerWithStats(routeName string, router http.Handler, logger Logger, statsd StatsD) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Debug(r.Method, "at", r.URL.String())
		metrics := httpsnoop.CaptureMetrics(router, w, r)

		logResult(routeName, metrics, statsd, logger, r)
	})
}

func logResult(routeName string, metrics httpsnoop.Metrics, statsd StatsD, logger Logger, req *http.Request) {
	responseTag := fmt.Sprintf("response:%d", metrics.Code)
	tags := []string{"route:" + routeName, responseTag}
	if caller := req.Header.Get("X-Component"); caller != "" {
		tags = append(tags, "caller:"+caller)
	}
	statsd.Histogram(WebResponseTimeKey, float64(metrics.Duration.Nanoseconds())/1000000, tags...)
	statsd.Incr(fmt.Sprintf(WebResponseCodeFormatKey, metrics.Code), tags...)
	statsd.Incr(WebResponseCodeAllKey, tags...)
	logger.Debugf("Request to %s had response code %d in %dms", req.URL.String(), metrics.Code, metrics.Duration.Milliseconds())
}
