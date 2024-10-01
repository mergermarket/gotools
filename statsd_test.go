package tools

import (
	"reflect"
	"testing"
)

func TestNewStatsD(t *testing.T) {
	config := StatsDConfig{
		isProduction: true,
		host:         "localhost",
		port:         "8080",
	}
	sd, _ := newMMStatsD(config)
	namespace := reflect.ValueOf(sd.ddstatsd).Elem().FieldByName("namespace").String()

	if namespace != "app." {
		t.Error("Was expecting the default namespace, but got ", namespace)
	}
}

func TestReturnsDummyStatsDIfConfigIsPoo(t *testing.T) {
	config := StatsDConfig{
		isProduction: true,
		host:         "",
		port:         "",
	}
	_, err := NewStatsD(config)

	if err != nil {
		t.Error("Was not expecting to get an error, yet error came")
	}
}

func TestGlobalTagging(t *testing.T) {
	config := StatsDConfig{
		isProduction: true,
		host:         "localhost",
		port:         "8080",
	}
	sd, _ := newMMStatsD(config)

	tags := reflect.ValueOf(sd.ddstatsd).Elem().FieldByName("tags").Slice(0, 2)

	if !contains(tags, "env:local") {
		t.Error("Expected a global tag of 'env:local', but it wasn't in ", tags)
	}

	if !contains(tags, "component:a-service-has-no-name") {
		t.Error("Expected a global tag of 'component:a-service-has-no-name', but it wasn't in ", tags)
	}
}

func contains(strings reflect.Value, target string) bool {
	for i := 0; i < strings.Len(); i++ {
		if strings.Index(i).String() == target {
			return true
		}
	}
	return false
}
