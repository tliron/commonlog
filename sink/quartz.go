package sink

import (
	"fmt"

	"github.com/reugn/go-quartz/logger"
	"github.com/tliron/commonlog"
	"github.com/tliron/kutil/util"
)

// Broken in Quartz v0.8.0, see: https://github.com/reugn/go-quartz/issues/74
func SetDefaultQuartzLogger(name ...string) {
	logger.SetDefault(NewQuartzLogger(name...))
}

//
// QuartzLogger
//

type QuartzLogger struct {
	name []string
}

func NewQuartzLogger(name ...string) *QuartzLogger {
	return &QuartzLogger{
		name: name,
	}
}

// ([logger.Logger] interface)
func (self *QuartzLogger) Trace(msg any) {
	self.sendMessage(commonlog.Debug, util.ToString(msg))
}

// ([logger.Logger] interface)
func (self *QuartzLogger) Tracef(format string, args ...any) {
	self.sendMessage(commonlog.Debug, fmt.Sprintf(format, args...))
}

// ([logger.Logger] interface)
func (self *QuartzLogger) Debug(msg any) {
	self.sendMessage(commonlog.Info, util.ToString(msg))
}

// ([logger.Logger] interface)
func (self *QuartzLogger) Debugf(format string, args ...any) {
	self.sendMessage(commonlog.Info, fmt.Sprintf(format, args...))
}

// ([logger.Logger] interface)
func (self *QuartzLogger) Info(msg any) {
	self.sendMessage(commonlog.Notice, util.ToString(msg))
}

// ([logger.Logger] interface)
func (self *QuartzLogger) Infof(format string, args ...any) {
	self.sendMessage(commonlog.Notice, fmt.Sprintf(format, args...))
}

// ([logger.Logger] interface)
func (self *QuartzLogger) Warn(msg any) {
	self.sendMessage(commonlog.Warning, util.ToString(msg))
}

// ([logger.Logger] interface)
func (self *QuartzLogger) Warnf(format string, args ...any) {
	self.sendMessage(commonlog.Warning, fmt.Sprintf(format, args...))
}

// ([logger.Logger] interface)
func (self *QuartzLogger) Error(msg any) {
	self.sendMessage(commonlog.Error, util.ToString(msg))
}

// ([logger.Logger] interface)
func (self *QuartzLogger) Errorf(format string, args ...any) {
	self.sendMessage(commonlog.Error, fmt.Sprintf(format, args...))
}

// ([logger.Logger] interface)
func (self *QuartzLogger) Enabled(level logger.Level) bool {
	return commonlog.AllowLevel(quartzToLevel(level), self.name...)
}

// Utils

func (self *QuartzLogger) sendMessage(level commonlog.Level, msg string) {
	if message := commonlog.NewMessage(level, 3, self.name...); message != nil {
		message.Set(commonlog.MESSAGE, msg)
		message.Send()
	}
}

func quartzToLevel(level logger.Level) commonlog.Level {
	switch level {
	case logger.LevelTrace:
		return commonlog.Debug
	case logger.LevelDebug:
		return commonlog.Info
	case logger.LevelInfo:
		return commonlog.Notice
	case logger.LevelWarn:
		return commonlog.Warning
	case logger.LevelError:
		return commonlog.Error
	case logger.LevelOff:
		return commonlog.None
	default:
		panic(fmt.Sprintf("unsupported log level: %d", level))
	}
}
