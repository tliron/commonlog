package commonlog

import (
	"github.com/tliron/kutil/util"
)

//
// Message
//

type Message interface {
	// Sets a value on the message and returns the same message
	// object.
	//
	// These keys are often specially supported:
	//
	// "message": the base text of the message
	// "scope": the scope of the message
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
	send    SendUnstructuredMessageFunc
}

func NewUnstructuredMessage(send SendUnstructuredMessageFunc) *UnstructuredMessage {
	return &UnstructuredMessage{
		send: send,
	}
}

// ([Message] interface)
func (self *UnstructuredMessage) Set(key string, value any) Message {
	switch key {
	case "message":
		self.message = util.ToString(value)

	case "scope":
		self.prefix = "{" + util.ToString(value) + "}"

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

	self.send(message)
}
