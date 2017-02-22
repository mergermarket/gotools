package tools

import (
	"errors"
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
		url := r.URL.String()
		logger.Debug(r.Method, "at", url)

		defer handlePanic(url, routeName, w, logger, statsd)

		metrics := httpsnoop.CaptureMetrics(router, w, r)

		logResult(routeName, metrics.Code, metrics.Duration.Nanoseconds(), statsd, logger, url)
	})
}

func handlePanic(url, routeName string, w http.ResponseWriter, logger logger, statsd StatsD) {
	var err error

	r := recover()
	if r != nil {
		switch t := r.(type) {
		case string:
			err = errors.New(t)
		case error:
			err = t
		default:
			err = errors.New("Unknown error")
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Error(fmt.Sprintf("unhandled panic serving %s %v", url, err))
		logResult(routeName, http.StatusInternalServerError, 0, statsd, logger, url)
	}
}

func logResult(routeName string, statusCode int, responseTime int64, statsd StatsD, logger logger, url string) {
	responseTag := fmt.Sprintf("response:%d", statusCode)
	statsd.Histogram(WebResponseTimeKey, float64(responseTime)/1000000, "route:"+routeName, responseTag)
	responseCodeKey := fmt.Sprintf(WebResponseCodeFormatKey, statusCode)
	statsd.Incr(responseCodeKey, "route:"+routeName, responseTag)
	statsd.Incr(WebResponseCodeAllKey, "route:"+routeName, responseTag)

	message := fmt.Sprint("Request to ", url, " had response code ", statusCode)
	if statusCode >= 400 {
		logger.Error(message)
	} else {
		logger.Info(message)
	}
}
