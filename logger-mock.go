package commonlog

//
// MockLogger
//

var MOCK_LOGGER MockLogger

// [Logger] that does nothing.
type MockLogger struct{}

// ([Logger] interface)
func (self MockLogger) AllowLevel(level Level) bool {
	return false
}

// ([Logger] interface)
func (self MockLogger) SetMaxLevel(level Level) {
}

// ([Logger] interface)
func (self MockLogger) GetMaxLevel() Level {
	return None
}

// ([Logger] interface)
func (self MockLogger) NewMessage(level Level, depth int, keysAndValues ...any) Message {
	return nil
}

// ([Logger] interface)
func (self MockLogger) Log(level Level, depth int, message string, keysAndValues ...any) {
}

// ([Logger] interface)
func (self MockLogger) Logf(level Level, depth int, format string, args ...any) {
}

// ([Logger] interface)
func (self MockLogger) Critical(message string, keysAndValues ...any) {
}

// ([Logger] interface)
func (self MockLogger) Criticalf(format string, args ...any) {
}

// ([Logger] interface)
func (self MockLogger) Error(message string, keysAndValues ...any) {
}

// ([Logger] interface)
func (self MockLogger) Errorf(format string, args ...any) {
}

// ([Logger] interface)
func (self MockLogger) Warning(message string, keysAndValues ...any) {
}

// ([Logger] interface)
func (self MockLogger) Warningf(format string, args ...any) {
}

// ([Logger] interface)
func (self MockLogger) Notice(message string, keysAndValues ...any) {
}

// ([Logger] interface)
func (self MockLogger) Noticef(format string, args ...any) {
}

// ([Logger] interface)
func (self MockLogger) Info(message string, keysAndValues ...any) {
}

// ([Logger] interface)
func (self MockLogger) Infof(format string, args ...any) {
}

// ([Logger] interface)
func (self MockLogger) Debug(message string, keysAndValues ...any) {
}

// ([Logger] interface)
func (self MockLogger) Debugf(format string, args ...any) {
}
