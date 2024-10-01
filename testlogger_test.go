package tools

import "testing"

func TestTestLogger(t *testing.T) {
	logger, _ := NewTestTools(t)
	t.Run("it prints things nicely", func(_ *testing.T) {
		logger.Debugf("hello %s and %s", "Baktash", "Olo")
	})
}
