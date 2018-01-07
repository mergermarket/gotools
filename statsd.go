package tools

import (
	"fmt"
	"os"

	"errors"
	"github.com/DataDog/datadog-go/statsd"
)

const (
	statsdRate                      = 1
	dummyFmtString1                 = "%s: name: %s, value: %f, tags: %v"
	dummyFmtString2                 = "%s: name: %s, tags: %v"
	dummyFmtString3                 = "%s: name: %s, value: %d, tags: %v"
	HttpClientResponseCodeAllKey    = "http_client.response_code.all"
	HttpClientResponseTimeKey       = "http_client.response_time_ms"
	HttpClientResponseErrorKey      = "http_client.response_error"
	HttpClientResponseSuccessKey    = "http_client.response_success"
	HttpClientResponseCodeFormatKey = "http_client.response_code.%d"
	WebResponseTimeKey              = "web.response_time"
	WebResponseCodeFormatKey        = "web.response_code.%d"
	WebResponseCodeAllKey           = "web.response_code.all"
)

// StatsD is the interface for the all the DataDog StatsD methods used by
// mergermarket. Please extend it as needed to record other types of metric.
type StatsD interface {
	Histogram(name string, value float64, tags ...string)
	Gauge(name string, value float64, tags ...string)
	Incr(name string, tags ...string)
	Count(name string, value int64, tags ...string)
}

// StatsDConfig provides configuration for metrics recording
type StatsDConfig struct {
	isProduction bool
	log          Logger
	host         string
	port         string
}

func NewStatsDConfig(isProduction bool, log Logger) StatsDConfig {
	return StatsDConfig{
		isProduction: isProduction,
		log:          log,
		host:         os.Getenv("STATSD_HOST"),
		port:         os.Getenv("STATSD_PORT"),
	}
}

// NewStatsD provides a new StatsD metrics recorder
func NewStatsD(config StatsDConfig) (StatsD, error) {
	if config.isProduction == false || config.port == "" || config.host == "" {
		return &dummyStatsD{config.log}, nil
	}
	return newMMStatsD(config)
}

func newMMStatsD(config StatsDConfig) (*mmStatsD, error) {
	if config.port == "" || config.host == "" {
		return nil, errors.New("Port and Host are required fields")
	}

	sd, err := statsd.New(config.host + ":" + config.port)

	if err != nil {
		return nil, err
	}

	addGlobalNamespace(sd)
	addGlobalTags(sd)

	return &mmStatsD{sd, config.log}, nil
}

func addGlobalNamespace(sd *statsd.Client) {
	sd.Namespace = "app."
}

func addGlobalTags(sd *statsd.Client) {
	sd.Tags = append(sd.Tags, globalTags()...)
}

func globalTags() []string {
	return []string{
		"env:" + getEnv(),
		"component:" + getComponentName(),
	}
}

type mmStatsD struct {
	ddstatsd *statsd.Client
	log      Logger
}

const statsDErrMsg = "Failed to send"
const statsDErrFmt = "%s %s %v"

func (mmsd *mmStatsD) Histogram(name string, value float64, tags ...string) {
	if err := mmsd.ddstatsd.Histogram(name, value, tags, statsdRate); err != nil {
		errMsg := fmt.Sprintf(statsDErrFmt, statsDErrMsg, name, err)
		mmsd.log.Error(errMsg)
	}
}

func (mmsd *mmStatsD) Gauge(name string, value float64, tags ...string) {
	if err := mmsd.ddstatsd.Gauge(name, value, tags, statsdRate); err != nil {
		errMsg := fmt.Sprintf(statsDErrFmt, statsDErrMsg, name, err)
		mmsd.log.Error(errMsg)
	}
}

func (mmsd *mmStatsD) Incr(name string, tags ...string) {
	if err := mmsd.ddstatsd.Incr(name, tags, statsdRate); err != nil {
		errMsg := fmt.Sprintf(statsDErrFmt, statsDErrMsg, name, err)
		mmsd.log.Error(errMsg)
	}
}

func (mmsd *mmStatsD) Count(name string, value int64, tags ...string) {
	if err := mmsd.ddstatsd.Count(name, value, tags, statsdRate); err != nil {
		errMsg := fmt.Sprintf(statsDErrFmt, statsDErrMsg, name, err)
		mmsd.log.Error(errMsg)
	}
}

// dummyStatsD is returned when StatsDConfig.isDevelopment is set to true. It
// stubs out the DataDog methods and sends them to the supplied logger
type dummyStatsD struct {
	Logger
}

func (dsd dummyStatsD) Histogram(name string, value float64, tags ...string) {
	logString := fmt.Sprintf(dummyFmtString1, "Histogram", name, value, tags)
	dsd.Info(logString)
}

func (dsd dummyStatsD) Gauge(name string, value float64, tags ...string) {
	logString := fmt.Sprintf(dummyFmtString1, "Gauge", name, value, tags)
	dsd.Info(logString)
}

func (dsd dummyStatsD) Incr(name string, tags ...string) {
	logString := fmt.Sprintf(dummyFmtString2, "Increment", name, tags)
	dsd.Info(logString)
}

func (dsd dummyStatsD) Count(name string, value int64, tags ...string) {
	logString := fmt.Sprintf(dummyFmtString3, "Count", name, value, tags)
	dsd.Info(logString)
}
