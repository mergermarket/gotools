package tools

import (
	"fmt"
	"os"

	"errors"
	"github.com/DataDog/datadog-go/statsd"
)

const statsdRate = 1
const dummyFmtString = "%s: name: %s, value: %f, tags: %v"

// StatsD is the interface for the all the DataDog StatsD methods used by
// mergermarket. Please extend it as needed to record other types of metric.
type StatsD interface {
	Histogram(name string, value float64, tags ...string) error
	Gauge(name string, value float64, tags ...string) error
	Incr(name string, tags ...string) error
}

// StatsDConfig provides configuration for metrics recording
type StatsDConfig struct {
	isProduction bool
	log          logger
	host         string
	port         string
}

func NewStatsDConfig(isProduction bool, log logger) StatsDConfig {
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
		return nil, errors.New("You bastard")
	}

	sd, err := statsd.New(config.host + ":" + config.port)

	if err != nil {
		return nil, err
	}

	addGlobalNamespace(sd)
	addGlobalTags(sd)

	return &mmStatsD{sd}, nil
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
}

func (mmsd *mmStatsD) Histogram(name string, value float64, tags ...string) error {
	return mmsd.ddstatsd.Histogram(name, value, tags, statsdRate)
}

func (mmsd *mmStatsD) Gauge(name string, value float64, tags ...string) error {
	return mmsd.ddstatsd.Gauge(name, value, tags, statsdRate)
}

func (mmsd *mmStatsD) Incr(name string, tags ...string) error {
	return mmsd.ddstatsd.Incr(name, tags, statsdRate)
}

// dummyStatsD is returned when StatsDConfig.isDevelopment is set to true. It
// stubs out the DataDog methods and sends them to the supplied logger
type dummyStatsD struct {
	logger
}

func (dsd dummyStatsD) Histogram(name string, value float64, tags ...string) error {
	logString := fmt.Sprintf(dummyFmtString, "Histogram", name, value, tags)
	dsd.Info(logString)
	return nil
}

func (dsd dummyStatsD) Gauge(name string, value float64, tags ...string) error {
	logString := fmt.Sprintf(dummyFmtString, "Gauge", name, value, tags)
	dsd.Info(logString)
	return nil
}

func (dsd dummyStatsD) Incr(name string, tags ...string) error {
	logString := fmt.Sprintf(dummyFmtString, "Increment", name, 0.0, tags)
	dsd.Info(logString)
	return nil
}
