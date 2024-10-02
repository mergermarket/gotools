package tools

import (
	"fmt"
	"net/http"
	"runtime"
	"time"
)

type goRoutines struct {
	statsd StatsD
}

// NewGoRoutines returns you a http handler to tell you how many go routines there are and reports every 5 seconds to statsd the number of routines you are running
func NewGoRoutines(statsd StatsD, componentName string) http.Handler {
	g := new(goRoutines)
	g.statsd = statsd

	statsDKey := fmt.Sprintf("%s.goroutines", componentName)

	go func() {
		for {
			time.Sleep(5 * time.Second)
			goroutines := runtime.NumGoroutine()
			g.statsd.Gauge(statsDKey, float64(goroutines))
		}
	}()

	return g
}

func (g *goRoutines) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	goroutines := runtime.NumGoroutine()
	fmt.Fprintf(w, "%d go routines", goroutines)
}
