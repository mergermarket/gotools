package tools

type T interface {
	Log(args ...interface{})
	Logf(format string, args ...interface{})
}

// TestLogger accepts the testing package so you wont be bombarded with logs
// when your tests pass but if they fail you will see what's going on.
type TestLogger struct {
	T T
}

// Info logs info to the test logger.
func (testLogger TestLogger) Info(msg ...interface{}) {
	testLogger.T.Log("[Info]", msg)
}

// Infof logs info to the test logger.
func (testLogger TestLogger) Infof(format string, msg ...interface{}) {
	testLogger.T.Logf("[Info] "+format, msg)
}

// Debug logs debug to the test logger.
func (testLogger TestLogger) Debug(msg ...interface{}) {
	testLogger.T.Log("[Debug]", msg)
}

// Debugf logs debug to the test logger.
func (testLogger TestLogger) Debugf(format string, msg ...interface{}) {
	testLogger.T.Log("[Debug] "+format, msg)
}

// Error logs error to the test logger.
func (testLogger TestLogger) Error(msg ...interface{}) {
	testLogger.T.Log("[Error]", msg)
}

// Errorf logs error to the test logger.
func (testLogger TestLogger) Errorf(format string, msg ...interface{}) {
	testLogger.T.Log("[Error] "+format, msg)
}

// Warn logs warn to the test logger.
func (testLogger TestLogger) Warn(msg ...interface{}) {
	testLogger.T.Log("[Warn]", msg)
}

// Warnf logs warn to the test logger.
func (testLogger TestLogger) Warnf(format string, msg ...interface{}) {
	testLogger.T.Log("[Warn] "+format, msg)
}
