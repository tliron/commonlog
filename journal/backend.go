package journal

import (
	"fmt"
	"io"
	"strings"

	"github.com/coreos/go-systemd/journal"
	"github.com/tliron/commonlog"
)

func init() {
	backend := NewBackend()
	backend.Configure(0, nil)
	commonlog.SetBackend(backend)
}

//
// Backend
//

type Backend struct {
	VarsInMessage bool

	nameHierarchy *commonlog.NameHierarchy
	writer        io.Writer
}

func NewBackend() *Backend {
	return &Backend{
		nameHierarchy: commonlog.NewNameHierarchy(),
	}
}

// ([commonlog.Backend] interface)
func (self *Backend) Configure(verbosity int, path *string) {
	maxLevel := commonlog.VerbosityToMaxLevel(verbosity)

	if maxLevel == commonlog.None {
		self.writer = io.Discard
		self.nameHierarchy.SetMaxLevel(commonlog.None)
	} else {
		self.writer = JournalWriter{}
		self.nameHierarchy.SetMaxLevel(maxLevel)
	}
}

// ([commonlog.Backend] interface)
func (self *Backend) GetWriter() io.Writer {
	return self.writer
}

// ([commonlog.Backend] interface)
func (self *Backend) NewMessage(level commonlog.Level, depth int, name ...string) commonlog.Message {
	if self.AllowLevel(level, name...) {
		var priority journal.Priority
		switch level {
		case commonlog.Critical:
			priority = journal.PriCrit
		case commonlog.Error:
			priority = journal.PriErr
		case commonlog.Warning:
			priority = journal.PriWarning
		case commonlog.Notice:
			priority = journal.PriNotice
		case commonlog.Info:
			priority = journal.PriInfo
		case commonlog.Debug:
			priority = journal.PriDebug
		default:
			panic(fmt.Sprintf("unsupported log level: %d", level))
		}

		var prefix string
		if name := strings.Join(name, "."); len(name) > 0 {
			prefix = "[" + name + "] "
		}

		return commonlog.TraceMessage(NewMessage(priority, prefix, self.VarsInMessage), depth)
	} else {
		return nil
	}

}

// ([commonlog.Backend] interface)
func (self *Backend) AllowLevel(level commonlog.Level, name ...string) bool {
	return self.nameHierarchy.AllowLevel(level, name...)
}

// ([commonlog.Backend] interface)
func (self *Backend) SetMaxLevel(level commonlog.Level, name ...string) {
	self.nameHierarchy.SetMaxLevel(level, name...)
}

// ([commonlog.Backend] interface)
func (self *Backend) GetMaxLevel(name ...string) commonlog.Level {
	return self.nameHierarchy.GetMaxLevel(name...)
}
