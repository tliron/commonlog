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
	priority journal.Priority
	prefix   string
	message  string
	vars     map[string]string
}

func NewMessage(priority journal.Priority, prefix string) commonlog.Message {
	return &Message{
		priority: priority,
		prefix:   prefix,
	}
}

// ([commonlog.Message] interface)
func (self *Message) Set(key string, value any) commonlog.Message {
	switch key {
	case commonlog.MESSAGE:
		self.message = util.ToString(value)

	default:
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
		self.vars[key] = util.ToString(value)
	}

	return self
}

// ([commonlog.Message] interface)
func (self *Message) Send() {
	journal.Send(self.prefix+self.message, self.priority, self.vars)
}
