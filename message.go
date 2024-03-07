package commonlog

import (
	"runtime"
	"strconv"
	"strings"

	"github.com/tliron/kutil/util"
)

const (
	MESSAGE = "_message"
	SCOPE   = "_scope"
	FILE    = "_file"
	LINE    = "_line"
)

func TraceMessage(message Message, depth int) Message {
	if Trace && (message != nil) {
		if _, file, line, ok := runtime.Caller(depth + 2); ok {
			message.Set(FILE, file)
			message.Set(LINE, line)
		}
	}

	return message
}

func SendMessageWithTrace(message Message, depth int) {
	if Trace {
		if _, file, line, ok := runtime.Caller(depth + 2); ok {
			message.Set(FILE, file)
			message.Set(LINE, line)
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
	case MESSAGE:
		self.Message = util.ToString(value)

	case SCOPE:
		self.Scope = util.ToString(value)

	case FILE:
		self.File = util.ToString(value)

	case LINE:
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
	var builder strings.Builder

	if len(self.Message) > 0 {
		builder.WriteString(self.Message)
	}

	if values := self.ValuesString(true); values != "" {
		if builder.Len() > 0 {
			builder.WriteRune(' ')
		}
		builder.WriteString(values)
	}

	return strings.ReplaceAll(builder.String(), "\n", "¶")
}

func (self *UnstructuredMessage) Prefix(name ...string) string {
	if (len(name) == 0) && (self.Scope == "") {
		return ""
	}

	var builder strings.Builder

	builder.WriteRune('[')

	switch length := len(name); length {
	case 0:
	case 1:
		builder.WriteString(name[0])
	default:
		last := length - 1
		for _, n := range name[:last] {
			builder.WriteString(n)
			builder.WriteRune('.')
		}
		builder.WriteString(name[last])
	}

	if self.Scope != "" {
		builder.WriteRune(':')
		builder.WriteString(self.Scope)
	}

	builder.WriteRune(']')

	return builder.String()
}

func (self *UnstructuredMessage) StringWithPrefix(name ...string) string {
	if prefix := self.Prefix(name...); prefix == "" {
		return self.String()
	} else {
		return prefix + " " + self.String()
	}
}

func (self *UnstructuredMessage) ValuesString(withLocation bool) string {
	if len(self.Values) == 0 {
		return ""
	}

	var values strings.Builder

	values.WriteRune('{')

	values_ := self.Values
	if withLocation {
		if self.File != "" {
			values_ = append(values_, UnstructuredValue{FILE, self.File})
			if self.Line != -1 {
				values_ = append(values_, UnstructuredValue{LINE, strconv.FormatInt(self.Line, 10)})
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

	values.WriteRune('}')

	return values.String()
}

func (self *UnstructuredMessage) LocationString() string {
	if self.File == "" {
		return ""
	}

	var location strings.Builder

	location.WriteString(self.File)
	if self.Line != -1 {
		location.WriteRune(':')
		location.WriteString(strconv.FormatInt(self.Line, 10))
	}

	return location.String()
}
