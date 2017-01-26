package tools

import (
	"math/rand"
	"testing"
	"time"
)

// NewTestTools returns a logger and statsd instance to use in your tests
func NewTestTools(t *testing.T) (Logger, StatsD) {
	logger := TestLogger{T: t}
	statsd, _ := NewStatsD(NewStatsDConfig(false, logger))
	return logger, statsd
}

// RandomString creates a random string of alphanumeric characters of length strlen
// from: https://siongui.github.io/2015/04/13/go-generate-random-string/
// there are faster ways: http://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang
func RandomString(strlen int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}
