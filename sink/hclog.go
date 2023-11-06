package sink

import (
	"fmt"
	"io"
	logpkg "log"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/tliron/commonlog"
)

//
// HCLogger
//

type HCLogger struct {
	name []string
	args []any
}

func NewHCLogger(args []any, name ...string) *HCLogger {
	return &HCLogger{
		name: name,
		args: args,
	}
}

// ([hclog.Logger] interface)
func (self *HCLogger) Log(level hclog.Level, msg string, args ...any) {
	self.sendMessage(level, msg, args)
}

// ([hclog.Logger] interface)
func (self *HCLogger) Trace(msg string, args ...any) {
	self.sendMessage(hclog.Trace, msg, args)
}

// ([hclog.Logger] interface)
func (self *HCLogger) Debug(msg string, args ...any) {
	self.sendMessage(hclog.Debug, msg, args)
}

// ([hclog.Logger] interface)
func (self *HCLogger) Info(msg string, args ...any) {
	self.sendMessage(hclog.Info, msg, args)
}

// ([hclog.Logger] interface)
func (self *HCLogger) Warn(msg string, args ...any) {
	self.sendMessage(hclog.Warn, msg, args)
}

// ([hclog.Logger] interface)
func (self *HCLogger) Error(msg string, args ...any) {
	self.sendMessage(hclog.Error, msg, args)
}

// ([hclog.Logger] interface)
func (self *HCLogger) IsTrace() bool {
	return commonlog.AllowLevel(commonlog.Debug, self.name...)
}

// ([hclog.Logger] interface)
func (self *HCLogger) IsDebug() bool {
	return commonlog.AllowLevel(commonlog.Info, self.name...)
}

// ([hclog.Logger] interface)
func (self *HCLogger) IsInfo() bool {
	return commonlog.AllowLevel(commonlog.Notice, self.name...)
}

// ([hclog.Logger] interface)
func (self *HCLogger) IsWarn() bool {
	return commonlog.AllowLevel(commonlog.Warning, self.name...)
}

// ([hclog.Logger] interface)
func (self *HCLogger) IsError() bool {
	return commonlog.AllowLevel(commonlog.Error, self.name...)
}

// ([hclog.Logger] interface)
func (self *HCLogger) ImpliedArgs() []any {
	return self.args
}

// ([hclog.Logger] interface)
func (self *HCLogger) With(args ...any) hclog.Logger {
	return NewHCLogger(args, self.name...)
}

// ([hclog.Logger] interface)
func (self *HCLogger) Name() string {
	return strings.Join(self.name, ".")
}

// ([hclog.Logger] interface)
func (self *HCLogger) Named(name string) hclog.Logger {
	return NewHCLogger(self.args, append(self.name, name)...)
}

// ([hclog.Logger] interface)
func (self *HCLogger) ResetNamed(name string) hclog.Logger {
	return NewHCLogger(self.args, name)
}

// ([hclog.Logger] interface)
func (self *HCLogger) SetLevel(level hclog.Level) {
	commonlog.SetMaxLevel(hcToLevel(level), self.name...)
}

// ([hclog.Logger] interface)
func (self *HCLogger) GetLevel() hclog.Level {
	return hcFromLevel(commonlog.GetMaxLevel(self.name...))
}

// ([hclog.Logger] interface)
func (self *HCLogger) StandardLogger(opts *hclog.StandardLoggerOptions) *logpkg.Logger {
	// TODO
	return nil
}

// ([hclog.Logger] interface)
func (self *HCLogger) StandardWriter(opts *hclog.StandardLoggerOptions) io.Writer {
	return commonlog.GetWriter()
}

// Utils

func (self *HCLogger) sendMessage(level hclog.Level, msg string, args []any) {
	if message := commonlog.NewMessage(hcToLevel(level), 2, self.name...); message != nil {
		message.Set("message", msg)

		args = append(self.args, args...)
		if length := len(args); length%2 == 0 {
			for i := 0; i < length; i += 2 {
				if key, ok := args[i].(string); ok {
					switch key {
					case "message", "scope":
					default:
						message.Set(key, args[i+1])
					}
				}
			}
		}

		message.Send()
	}
}

func hcToLevel(level hclog.Level) commonlog.Level {
	switch level {
	case hclog.NoLevel:
		return commonlog.None
	case hclog.Trace:
		return commonlog.Debug
	case hclog.Debug:
		return commonlog.Info
	case hclog.Info:
		return commonlog.Notice
	case hclog.Warn:
		return commonlog.Warning
	case hclog.Error:
		return commonlog.Error
	default:
		panic(fmt.Sprintf("unsupported level: %d", level))
	}
}

func hcFromLevel(level commonlog.Level) hclog.Level {
	switch level {
	case commonlog.None:
		return hclog.NoLevel
	case commonlog.Critical:
		return hclog.Error
	case commonlog.Error:
		return hclog.Error
	case commonlog.Warning:
		return hclog.Warn
	case commonlog.Notice:
		return hclog.Info
	case commonlog.Info:
		return hclog.Debug
	case commonlog.Debug:
		return hclog.Trace
	default:
		panic(fmt.Sprintf("unsupported level: %d", level))
	}
}
