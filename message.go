package commonlog

import (
	"runtime"
	"strconv"
	"strings"

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

type SendUnstructuredMessageFunc func(message *UnstructuredMessage)

// Convenience type for implementing unstructured backends. Converts a structured
// message to an unstructured string.
type UnstructuredMessage struct {
	Scope   string
	Message string
	Values  []UnstructuredValue
	File    string
	Line    int64

	send SendUnstructuredMessageFunc
}

type UnstructuredValue struct {
	Key   string
	Value string
}

func NewUnstructuredMessage(send SendUnstructuredMessageFunc) *UnstructuredMessage {
	return &UnstructuredMessage{
		send: send,
		Line: -1,
	}
}

// ([Message] interface)
func (self *UnstructuredMessage) Set(key string, value any) Message {
	switch key {
	case "_message":
		self.Message = util.ToString(value)

	case "_scope":
		self.Scope = util.ToString(value)

	case "_file":
		self.File = util.ToString(value)

	case "_line":
		self.Line, _ = util.ToInt64(value)

	default:
		self.Values = append(self.Values, UnstructuredValue{key, util.ToString(value)})
	}

	return self
}

// ([Message] interface)
func (self *UnstructuredMessage) Send() {
	self.send(self)
}

// ([fmt.Stringify] interface)
func (self *UnstructuredMessage) String() string {
	var s strings.Builder

	if scope := self.ScopeString(); scope != "" {
		s.WriteString(scope)
	}

	if len(self.Message) > 0 {
		if s.Len() > 0 {
			s.WriteRune(' ')
		}
		s.WriteString(self.Message)
	}

	if values := self.ValuesString(true); values != "" {
		if s.Len() > 0 {
			s.WriteString("; ")
		}
		s.WriteString(values)
	}

	return strings.ReplaceAll(s.String(), "\n", "¶")
}

func (self *UnstructuredMessage) StringWithName(name ...string) string {
	s := self.String()
	if len(name) > 0 {
		s = "[" + strings.Join(name, ".") + "] " + s
	}
	return s
}

func (self *UnstructuredMessage) ScopeString() string {
	if len(self.Scope) > 0 {
		return "{" + self.Scope + "}"
	} else {
		return ""
	}
}

func (self *UnstructuredMessage) ValuesString(withLocation bool) string {
	var values strings.Builder

	values_ := self.Values
	if withLocation {
		if self.File != "" {
			values_ = append(values_, UnstructuredValue{"_file", self.File})
			if self.Line != -1 {
				values_ = append(values_, UnstructuredValue{"_line", strconv.FormatInt(self.Line, 10)})
			}
		}
	}

	last := len(values_) - 1
	for index, value := range values_ {
		values.WriteString(value.Key)
		values.WriteRune('=')
		values.WriteString(strconv.Quote(value.Value))
		if index != last {
			values.WriteRune(' ')
		}
	}

	return values.String()
}

func (self *UnstructuredMessage) LocationString() string {
	var location strings.Builder

	if self.File != "" {
		location.WriteString(self.File)
		if self.Line != -1 {
			location.WriteRune(':')
			location.WriteString(strconv.FormatInt(self.Line, 10))
		}
	}

	return location.String()
}
