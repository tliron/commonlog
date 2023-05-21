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
	hierarchy *commonlog.Hierarchy
	writer    io.Writer
}

func NewBackend() *Backend {
	return &Backend{
		hierarchy: commonlog.NewMaxLevelHierarchy(),
	}
}

// commonlog.Backend interface
func (self *Backend) Configure(verbosity int, path *string) {
	maxLevel := commonlog.VerbosityToMaxLevel(verbosity)

	if maxLevel == commonlog.None {
		self.writer = io.Discard
		self.hierarchy.SetMaxLevel(nil, commonlog.None)
	} else {
		self.writer = JournalWriter{}
		self.hierarchy.SetMaxLevel(nil, maxLevel)
	}
}

// commonlog.Backend interface
func (self *Backend) GetWriter() io.Writer {
	return self.writer
}

// commonlog.Backend interface
func (self *Backend) NewMessage(name []string, level commonlog.Level, depth int) commonlog.Message {
	if self.AllowLevel(name, level) {
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
			panic(fmt.Sprintf("unsupported level: %d", level))
		}

		var prefix string
		if name := strings.Join(name, "."); len(name) > 0 {
			prefix = "[" + name + "] "
		}

		return NewMessage(priority, prefix)
	} else {
		return nil
	}

}

// commonlog.Backend interface
func (self *Backend) AllowLevel(name []string, level commonlog.Level) bool {
	return self.hierarchy.AllowLevel(name, level)
}

// commonlog.Backend interface
func (self *Backend) SetMaxLevel(name []string, level commonlog.Level) {
	self.hierarchy.SetMaxLevel(name, level)
}

// commonlog.Backend interface
func (self *Backend) GetMaxLevel(name []string) commonlog.Level {
	return self.hierarchy.GetMaxLevel(name)
}
