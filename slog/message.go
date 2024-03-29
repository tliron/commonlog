package slog

import (
	contextpkg "context"
	"log/slog"

	"github.com/tliron/commonlog"
	"github.com/tliron/kutil/util"
)

//
// Message
//

type Message struct {
	logger  *slog.Logger
	level   slog.Level
	context contextpkg.Context

	message string
	args    []any
}

func NewMessage(logger *slog.Logger, level slog.Level, context contextpkg.Context) commonlog.Message {
	return &Message{
		logger:  logger,
		level:   level,
		context: context,
	}
}

// ([commonlog.Message] interface)
func (self *Message) Set(key string, value any) commonlog.Message {
	switch key {
	case commonlog.MESSAGE:
		self.message = util.ToString(value)

	default:
		self.args = append(self.args, key, value)
	}

	return self
}

// ([commonlog.Message] interface)
func (self *Message) Send() {
	self.logger.Log(self.context, self.level, self.message, self.args...)
}
