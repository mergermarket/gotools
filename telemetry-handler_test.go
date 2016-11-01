package tools

import (
	"testing"
	"time"
)

func TestTimeTrack(t *testing.T) {
	msd := &MockStatsD{}
	tl := TestLogger{t}
	timeTrack("aTestRoute", time.Now(), tl, msd)

	call, err := msd.Call()

	if err != nil {
		t.Error("No call made to MockStatsD")
	}

	if call.Method != "Histogram" {
		t.Error("Expected a call to Histogram")
	}

	if call.Args.Name != "web.response_time" {
		t.Error("Expected name of metric to be 'web.response_time', but got:", call.Args.Name)
	}
	if !contains(call.Args.Tags, "route:aTestRoute") {
		t.Error("Expected to get a tag 'route:aTestRoute', but got:", call.Args.Tags)
	}
}
