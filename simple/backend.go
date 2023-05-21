package simple

import (
	"io"
	"os"

	"github.com/tliron/commonlog"
	"github.com/tliron/kutil/terminal"
	"github.com/tliron/kutil/util"
)

const LOG_FILE_WRITE_PERMISSIONS = 0600

const BUFFER_SIZE = 10000

func init() {
	backend := NewBackend()
	backend.Configure(0, nil)
	commonlog.SetBackend(backend)
}

//
// Backend
//

type Backend struct {
	Writer   io.Writer
	Format   FormatFunc
	Buffered bool

	colorize  bool
	hierarchy *commonlog.Hierarchy
}

func NewBackend() *Backend {
	return &Backend{
		Format:    DefaultFormat,
		Buffered:  true,
		hierarchy: commonlog.NewMaxLevelHierarchy(),
	}
}

// commonlog.Backend interface
func (self *Backend) Configure(verbosity int, path *string) {
	maxLevel := commonlog.VerbosityToMaxLevel(verbosity)

	if maxLevel == commonlog.None {
		self.Writer = io.Discard
		self.hierarchy.SetMaxLevel(nil, commonlog.None)
	} else {
		if path != nil {
			if file, err := os.OpenFile(*path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, LOG_FILE_WRITE_PERMISSIONS); err == nil {
				util.OnExitError(file.Close)
				if self.Buffered {
					writer := util.NewBufferedWriter(file, BUFFER_SIZE)
					util.OnExitError(writer.Close)
					self.Writer = writer
				} else {
					self.Writer = util.NewSyncedWriter(file)
				}
			} else {
				util.Failf("log file error: %s", err.Error())
			}
		} else {
			self.colorize = terminal.Colorize
			if self.Buffered {
				writer := util.NewBufferedWriter(os.Stderr, BUFFER_SIZE)
				util.OnExitError(writer.Close)
				self.Writer = writer
			} else {
				self.Writer = util.NewSyncedWriter(os.Stderr)
			}
		}

		self.hierarchy.SetMaxLevel(nil, maxLevel)
	}
}

// commonlog.Backend interface
func (self *Backend) GetWriter() io.Writer {
	return self.Writer
}

// commonlog.Backend interface
func (self *Backend) NewMessage(name []string, level commonlog.Level, depth int) commonlog.Message {
	if self.AllowLevel(name, level) {
		return commonlog.NewUnstructuredMessage(func(message string) {
			message = self.Format(message, name, level, self.colorize)
			io.WriteString(self.Writer, message+"\n")
		})
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
