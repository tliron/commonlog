package klog

import (
	"fmt"
	"io"
	"os"

	"github.com/tliron/commonlog"
	"github.com/tliron/go-kutil/util"
	"k8s.io/klog/v2"
)

const (
	LogFileWritePermissions = 0600
	DefaultBufferSize       = 1_000
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
	BufferSize int
	Buffered   bool

	writer        io.Writer
	nameHierarchy *commonlog.NameHierarchy
}

func NewBackend() *Backend {
	return &Backend{
		BufferSize:    DefaultBufferSize,
		Buffered:      true,
		nameHierarchy: commonlog.NewNameHierarchy(),
	}
}

var flushHandle util.ExitFunctionHandle

// ([commonlog.Backend] interface)
func (self *Backend) Configure(verbosity int, path *string) {
	// klog can also do its own configuration via klog.InitFlags

	if flushHandle == 0 {
		flushHandle = util.OnExit(klog.Flush)
	}

	maxLevel := commonlog.VerbosityToMaxLevel(verbosity)

	if maxLevel == commonlog.None {
		klog.SetOutput(io.Discard)
		self.nameHierarchy.SetMaxLevel(commonlog.None)
	} else {
		if path != nil {
			if file, err := os.OpenFile(*path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, LogFileWritePermissions); err == nil {
				util.OnExitError(file.Close)
				if self.Buffered {
					writer := util.NewBufferedWriter(file, self.BufferSize, false)
					util.OnExitError(writer.Close)
					self.writer = writer
					klog.SetOutput(writer)
				} else {
					klog.SetOutput(util.NewSyncedWriter(file))
				}
			} else {
				util.Failf("log file error: %s", err.Error())
			}
		} else if self.Buffered {
			writer := util.NewBufferedWriter(os.Stderr, self.BufferSize, false)
			util.OnExitError(writer.Close)
			self.writer = writer
			klog.SetOutput(writer)
		} else {
			klog.SetOutput(util.NewSyncedWriter(os.Stderr))
		}

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
		return commonlog.TraceMessage(commonlog.NewLinearMessage(func(message *commonlog.LinearMessage) {
			message_ := message.StringWithPrefix(name...)

			switch level {
			case commonlog.Critical:
				klog.ErrorDepth(depth, message_)
			case commonlog.Error:
				klog.ErrorDepth(depth, message_)
			case commonlog.Warning:
				klog.WarningDepth(depth, message_)
			case commonlog.Notice:
				klog.InfoDepth(depth, message_)
			case commonlog.Info:
				klog.InfoDepth(depth, message_)
			case commonlog.Debug:
				klog.InfoDepth(depth, message_)
			default:
				panic(fmt.Sprintf("unsupported log level: %d", level))
			}
		}), depth)
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
