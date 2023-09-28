package sink

import (
	contextpkg "context"
	"fmt"
	"log/slog"

	"github.com/tliron/commonlog"
)

func NewStandardStructuredLogger(name ...string) *slog.Logger {
	return slog.New(NewStandardStructuredHandler(name...))
}

//
// StandardStructuredHandler
//

type StandardStructuredHandler struct {
	name []string

	attrs             []slog.Attr
	currentGroupName  string
	currentGroupAttrs []slog.Attr
}

func NewStandardStructuredHandler(name ...string) *StandardStructuredHandler {
	return &StandardStructuredHandler{
		name: name,
	}
}

// ([slog.Handler] interface)
func (self *StandardStructuredHandler) Enabled(context contextpkg.Context, level slog.Level) bool {
	return commonlog.AllowLevel(slogToLevel(level), self.name...)
}

// ([slog.Handler] interface)
func (self *StandardStructuredHandler) Handle(context contextpkg.Context, record slog.Record) error {
	if message := commonlog.NewMessage(slogToLevel(record.Level), 2, self.name...); message != nil {
		message.Set("message", record.Message)

		self.resolve(false)
		for _, attr := range self.attrs {
			slogSet(message, "", attr)
		}

		message.Send()
	}

	return nil
}

// ([slog.Handler] interface)
func (self *StandardStructuredHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return self
	}

	self = self.clone()
	if self.currentGroupName == "" {
		self.attrs = append(self.attrs, attrs...)
	} else {
		self.currentGroupAttrs = append(self.currentGroupAttrs, attrs...)
	}

	return self
}

// ([slog.Handler] interface)
func (self *StandardStructuredHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return self
	}

	self = self.clone()
	self.resolve(true)
	self.currentGroupName = name

	return self
}

func (self *StandardStructuredHandler) clone() *StandardStructuredHandler {
	handler := NewStandardStructuredHandler(self.name...)

	handler.attrs = append(self.attrs[:0:0], self.attrs...)

	handler.currentGroupName = self.currentGroupName
	handler.currentGroupAttrs = append(self.currentGroupAttrs[:0:0], self.currentGroupAttrs...)

	return handler
}

func (self *StandardStructuredHandler) resolve(closeGroup bool) {
	if self.currentGroupName != "" {
		if len(self.currentGroupAttrs) > 0 {
			self.attrs = append(self.attrs, slog.Attr{
				Key:   self.currentGroupName,
				Value: slog.GroupValue(self.currentGroupAttrs...),
			})
			self.currentGroupAttrs = nil
		}

		if closeGroup {
			self.currentGroupName = ""
		}
	}
}

// Utils

func slogToLevel(level slog.Level) commonlog.Level {
	switch level {
	case slog.LevelDebug:
		return commonlog.Debug
	case slog.LevelInfo:
		return commonlog.Info
	case slog.LevelWarn:
		return commonlog.Warning
	case slog.LevelError:
		return commonlog.Error
	default:
		panic(fmt.Sprintf("unsupported level: %d", level))
	}
}

func slogSet(message commonlog.Message, prefix string, attr slog.Attr) {
	switch attr.Value.Kind() {
	case slog.KindGroup:
		if prefix != "" {
			prefix += "." + prefix
		}

		for _, attr_ := range attr.Value.Group() {
			slogSet(message, prefix, attr_)
		}

	default:
		value := attr.Value.Resolve().Any()
		if prefix == "" {
			message.Set(attr.Key, value)
		} else {
			message.Set(prefix+attr.Key, value)
		}
	}
}
