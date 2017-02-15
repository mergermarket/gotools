package tools

import (
	"testing"
	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestLogger_withFileAndLine(t *testing.T) {
	fields := logrus.Fields{}
 	updatedFields := withFileAndLine(fields)
	assert.NotEmpty(t, updatedFields["file"])
	assert.NotEmpty(t, updatedFields["line"])
}
