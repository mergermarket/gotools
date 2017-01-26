package testtools

import (
	"github.com/mergermarket/gotools/logging"
	"github.com/mergermarket/gotools/statsd"
	"testing"
)

func NewTestTools(t *testing.T) (logging.Logger, statsd.StatsD) {
	logger := TestLogger{T: t}
	statsd, _ := statsd.NewStatsD(statsd.NewStatsDConfig(false, logger))
	return logger, statsd
}
