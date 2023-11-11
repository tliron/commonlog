package commonlog

import (
	"runtime"
	"strconv"

	"github.com/tliron/kutil/util"
)

func TraceMessage(message Message, depth int) Message {
	if Trace && (message != nil) {
		if _, file, line, ok := runtime.Caller(depth + 2); ok {
			message.Set("_file", file)
			message.Set("_line", line)
		}
	}

	return message
}

func SendMessageWithTrace(message Message, depth int) {
	if Trace {
		if _, file, line, ok := runtime.Caller(depth + 2); ok {
			message.Set("_file", file)
			message.Set("_line", line)
		}
	}

	message.Send()
}

//
// Message
//

type Message interface {
	// Sets a value on the message and returns the same message
	// object.
	//
	// These keys are often specially supported:
	//
	// "_message": the base text of the message
	// "_scope": the scope of the message
	Set(key string, value any) Message

	// Sends the message to the backend
	Send()
}

//
// UnstructuredMessage
//

type SendUnstructuredMessageFunc func(message string)

// Convenience type for implementing unstructured backends. Converts a structured
// message to an unstructured string.
type UnstructuredMessage struct {
	prefix  string
	message string
	suffix  string
	file    string
	line    int64
	send    SendUnstructuredMessageFunc
}

func NewUnstructuredMessage(send SendUnstructuredMessageFunc) *UnstructuredMessage {
	return &UnstructuredMessage{
		send: send,
		line: -1,
	}
}

// ([Message] interface)
func (self *UnstructuredMessage) Set(key string, value any) Message {
	switch key {
	case "_message":
		self.message = util.ToString(value)

	case "_scope":
		self.prefix = "{" + util.ToString(value) + "}"

	case "_file":
		self.file = util.ToString(value)

	case "_line":
		self.line, _ = util.ToInt64(value)

	default:
		if len(self.suffix) > 0 {
			self.suffix += ", "
		}
		self.suffix += key + "=" + util.ToString(value)
	}

	return self
}

// ([Message] interface)
func (self *UnstructuredMessage) Send() {
	message := self.prefix

	if len(self.message) > 0 {
		if len(message) > 0 {
			message += " "
		}
		message += self.message
	}

	if len(self.suffix) > 0 {
		if len(message) > 0 {
			message += " "
		}
		message += self.suffix
	}

	if self.file != "" {
		message += "\n└─" + self.file
		if self.line != -1 {
			message += ":" + strconv.FormatInt(self.line, 10)
		}
	}

	self.send(message)
}
