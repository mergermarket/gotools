package tools

import (
	"testing"
)

func NewTestTools(t *testing.T) (Logger, StatsD) {
	logger := TestLogger{T: t}
	statsd, _ := NewStatsD(NewStatsDConfig(false, logger))
	return logger, statsd
}
