package commonlog

import (
	"fmt"

	"github.com/tliron/kutil/util"
)

//
// Logger
//

// While [NewMessage] is the "true" API entry point, this interface enables
// a familiar logger API. Because it's an interface, references can easily
// replace the implementation, for example setting a reference to
// [MOCK_LOGGER] will disable the logger.
//
// See [GetLogger].
type Logger interface {
	// Returns true if a level is loggable for this logger.
	AllowLevel(level Level) bool

	// Sets the maximum loggable level for this logger.
	SetMaxLevel(level Level)

	// Gets the maximum loggable level for this logger.
	GetMaxLevel() Level

	// Creates a new message for this logger. Will return nil if
	// the level is not loggable.
	//
	// The depth argument is used for skipping frames in callstack
	// logging, if supported.
	NewMessage(level Level, depth int, keysAndValues ...any) Message

	// Convenience method to create and send a message with at least
	// the "_message" key. Additional keys can be set by providing
	// a sequence of key-value pairs.
	Log(level Level, depth int, message string, keysAndValues ...any)

	// Convenience method to create and send a message with just
	// the "_message" key, where the message is created via the format
	// and args similarly to fmt.Printf.
	Logf(level Level, depth int, format string, args ...any)

	Critical(message string, keysAndValues ...any)
	Criticalf(format string, args ...any)
	Error(message string, keysAndValues ...any)
	Errorf(format string, args ...any)
	Warning(message string, keysAndValues ...any)
	Warningf(format string, args ...any)
	Notice(message string, keysAndValues ...any)
	Noticef(format string, args ...any)
	Info(message string, keysAndValues ...any)
	Infof(format string, args ...any)
	Debug(message string, keysAndValues ...any)
	Debugf(format string, args ...any)
}

//
// BackendLogger
//

// Default [Logger] implementation that logs to the current backend set with
// [SetBackend].
type BackendLogger struct {
	name []string
}

func NewBackendLogger(name ...string) BackendLogger {
	return BackendLogger{name: name}
}

// ([Logger] interface)
func (self BackendLogger) AllowLevel(level Level) bool {
	return AllowLevel(level, self.name...)
}

// ([Logger] interface)
func (self BackendLogger) SetMaxLevel(level Level) {
	SetMaxLevel(level, self.name...)
}

// ([Logger] interface)
func (self BackendLogger) GetMaxLevel() Level {
	return GetMaxLevel(self.name...)
}

// ([Logger] interface)
func (self BackendLogger) NewMessage(level Level, depth int, keysAndValues ...any) Message {
	if message := NewMessage(level, depth+1, self.name...); message != nil {
		SetMessageKeysAndValues(message, keysAndValues...)
		return message
	} else {
		return nil
	}
}

// ([Logger] interface)
func (self BackendLogger) Log(level Level, depth int, message string, keysAndValues ...any) {
	if message_ := self.NewMessage(level, depth+1, keysAndValues...); message_ != nil {
		message_.Set("_message", message)
		message_.Send()
	}
}

// ([Logger] interface)
func (self BackendLogger) Logf(level Level, depth int, format string, args ...any) {
	if message := self.NewMessage(level, depth+1); message != nil {
		message.Set("_message", fmt.Sprintf(format, args...))
		message.Send()
	}
}

// ([Logger] interface)
func (self BackendLogger) Critical(message string, keysAndValues ...any) {
	self.Log(Critical, 1, message, keysAndValues...)
}

// ([Logger] interface)
func (self BackendLogger) Criticalf(format string, args ...any) {
	self.Logf(Critical, 1, format, args...)
}

// ([Logger] interface)
func (self BackendLogger) Error(message string, keysAndValues ...any) {
	self.Log(Error, 1, message, keysAndValues...)
}

// ([Logger] interface)
func (self BackendLogger) Errorf(format string, args ...any) {
	self.Logf(Error, 1, format, args...)
}

// ([Logger] interface)
func (self BackendLogger) Warning(message string, keysAndValues ...any) {
	self.Log(Warning, 1, message, keysAndValues...)
}

// ([Logger] interface)
func (self BackendLogger) Warningf(format string, args ...any) {
	self.Logf(Warning, 1, format, args...)
}

// ([Logger] interface)
func (self BackendLogger) Notice(message string, keysAndValues ...any) {
	self.Log(Notice, 1, message, keysAndValues...)
}

// ([Logger] interface)
func (self BackendLogger) Noticef(format string, args ...any) {
	self.Logf(Notice, 1, format, args...)
}

// ([Logger] interface)
func (self BackendLogger) Info(message string, keysAndValues ...any) {
	self.Log(Info, 1, message, keysAndValues...)
}

// ([Logger] interface)
func (self BackendLogger) Infof(format string, args ...any) {
	self.Logf(Info, 1, format, args...)
}

// ([Logger] interface)
func (self BackendLogger) Debug(message string, keysAndValues ...any) {
	self.Log(Debug, 1, message, keysAndValues...)
}

// ([Logger] interface)
func (self BackendLogger) Debugf(format string, args ...any) {
	self.Logf(Debug, 1, format, args...)
}

//
// KeyValueLogger
//

type KeyValueLogger struct {
	logger        Logger
	keysAndValues []any
}

// Wrapping [Logger] that calls [Message.Set] with keys and values
// on all messages.
//
// If we're wrapping another [KeyValueLogger] then our keys and values
// will be merged into the wrapped keys and values.
func NewKeyValueLogger(logger Logger, keysAndValues ...any) KeyValueLogger {
	if keyValueLogger, ok := logger.(KeyValueLogger); ok {
		logger = keyValueLogger.logger
		keysAndValues, _ = MergeKeysAndValues(keyValueLogger.keysAndValues, keysAndValues)
	}

	return KeyValueLogger{
		logger:        logger,
		keysAndValues: keysAndValues,
	}
}

// Wrapping [Logger] that calls [Message.Set] with a "_scope" key
// on all messages.
//
// If we're wrapping another [KeyValueLogger] then our "_scope" key
// will be appended to the wrapped "_scope" key with a "." notation.
func NewScopeLogger(logger Logger, scope string) KeyValueLogger {
	var keysAndValues []any

	if keyValueLogger, ok := logger.(KeyValueLogger); ok {
		if scope_, ok := GetKeysAndValues(keyValueLogger.keysAndValues, "_scope"); ok {
			if scope_ != "" {
				scope = util.ToString(scope_) + "." + scope
			}
		}

		logger = keyValueLogger.logger
		keysAndValues, _ = MergeKeysAndValues(keyValueLogger.keysAndValues, []any{"_scope", scope})
	} else {
		keysAndValues = []any{"_scope", scope}
	}

	return KeyValueLogger{
		logger:        logger,
		keysAndValues: keysAndValues,
	}
}

// ([Logger] interface)
func (self KeyValueLogger) AllowLevel(level Level) bool {
	return self.logger.AllowLevel(level)
}

// ([Logger] interface)
func (self KeyValueLogger) SetMaxLevel(level Level) {
	self.logger.SetMaxLevel(level)
}

// ([Logger] interface)
func (self KeyValueLogger) GetMaxLevel() Level {
	return self.logger.GetMaxLevel()
}

// ([Logger] interface)
func (self KeyValueLogger) NewMessage(level Level, depth int, keysAndValues ...any) Message {
	if message := self.logger.NewMessage(level, depth+1, keysAndValues...); message != nil {
		SetMessageKeysAndValues(message, self.keysAndValues...)
		return message
	} else {
		return nil
	}
}

// ([Logger] interface)
func (self KeyValueLogger) Log(level Level, depth int, message string, keysAndValues ...any) {
	if message_ := self.NewMessage(level, depth+1, keysAndValues...); message_ != nil {
		message_.Set("_message", message)
		message_.Send()
	}
}

// ([Logger] interface)
func (self KeyValueLogger) Logf(level Level, depth int, format string, args ...any) {
	if message := self.NewMessage(level, depth+1); message != nil {
		message.Set("_message", fmt.Sprintf(format, args...))
		message.Send()
	}
}

// ([Logger] interface)
func (self KeyValueLogger) Critical(message string, keysAndValues ...any) {
	self.Log(Critical, 1, message, keysAndValues...)
}

// ([Logger] interface)
func (self KeyValueLogger) Criticalf(format string, args ...any) {
	self.Logf(Critical, 1, format, args...)
}

// ([Logger] interface)
func (self KeyValueLogger) Error(message string, keysAndValues ...any) {
	self.Log(Error, 1, message, keysAndValues...)
}

// ([Logger] interface)
func (self KeyValueLogger) Errorf(format string, args ...any) {
	self.Logf(Error, 1, format, args...)
}

// ([Logger] interface)
func (self KeyValueLogger) Warning(message string, keysAndValues ...any) {
	self.Log(Warning, 1, message, keysAndValues...)
}

// ([Logger] interface)
func (self KeyValueLogger) Warningf(format string, args ...any) {
	self.Logf(Warning, 1, format, args...)
}

// ([Logger] interface)
func (self KeyValueLogger) Notice(message string, keysAndValues ...any) {
	self.Log(Notice, 1, message, keysAndValues...)
}

// ([Logger] interface)
func (self KeyValueLogger) Noticef(format string, args ...any) {
	self.Logf(Notice, 1, format, args...)
}

// ([Logger] interface)
func (self KeyValueLogger) Info(message string, keysAndValues ...any) {
	self.Log(Info, 1, message, keysAndValues...)
}

// ([Logger] interface)
func (self KeyValueLogger) Infof(format string, args ...any) {
	self.Logf(Info, 1, format, args...)
}

// ([Logger] interface)
func (self KeyValueLogger) Debug(message string, keysAndValues ...any) {
	self.Log(Debug, 1, message, keysAndValues...)
}

// ([Logger] interface)
func (self KeyValueLogger) Debugf(format string, args ...any) {
	self.Logf(Debug, 1, format, args...)
}

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
