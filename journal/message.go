package journal

import (
	"strings"

	"github.com/coreos/go-systemd/journal"
	"github.com/tliron/commonlog"
	"github.com/tliron/kutil/util"
)

//
// Message
//

type Message struct {
	priority      journal.Priority
	prefix        string
	postfix       string
	message       string
	vars          map[string]string
	varsInMessage bool
}

func NewMessage(priority journal.Priority, prefix string, varsInMessage bool) commonlog.Message {
	return &Message{
		priority:      priority,
		prefix:        prefix,
		varsInMessage: varsInMessage,
	}
}

// ([commonlog.Message] interface)
func (self *Message) Set(key string, value any) commonlog.Message {
	value_ := util.ToString(value)

	switch key {
	case commonlog.MESSAGE:
		self.message = value_

	default:
		if self.varsInMessage {
			if self.postfix != "" {
				self.postfix += " "
			}
			self.postfix += key + "=" + value_
		}

		// See: https://www.freedesktop.org/software/systemd/man/systemd.journal-fields.html
		switch key {
		case commonlog.FILE:
			key = "CODE_FILE"

		case commonlog.LINE:
			key = "CODE_LINE"

		default:
			key = strings.ToUpper(key)
			if strings.HasPrefix(key, "_") {
				// Field name cannot be set by user
				key = "X" + key
			}
		}

		if self.vars == nil {
			self.vars = make(map[string]string)
		}
		self.vars[key] = value_
	}

	return self
}

// ([commonlog.Message] interface)
func (self *Message) Send() {
	message := self.prefix + self.message
	if self.postfix != "" {
		if message != "" {
			message += " "
		}
		message += "{" + self.postfix + "}"
	}
	journal.Send(message, self.priority, self.vars)
}
