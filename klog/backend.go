package klog

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/tliron/commonlog"
	"github.com/tliron/kutil/util"
	"k8s.io/klog/v2"
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
	Buffered bool

	writer    io.Writer
	hierarchy *commonlog.Hierarchy
}

func NewBackend() *Backend {
	return &Backend{
		Buffered:  true,
		hierarchy: commonlog.NewMaxLevelHierarchy(),
	}
}

var flushHandle util.ExitFunctionHandle

// commonlog.Backend interface
func (self *Backend) Configure(verbosity int, path *string) {
	// klog can also do its own configuration via klog.InitFlags

	if flushHandle == 0 {
		flushHandle = util.OnExit(klog.Flush)
	}

	maxLevel := commonlog.VerbosityToMaxLevel(verbosity)

	if maxLevel == commonlog.None {
		klog.SetOutput(io.Discard)
		self.hierarchy.SetMaxLevel(nil, commonlog.None)
	} else {
		if path != nil {
			if file, err := os.OpenFile(*path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, LOG_FILE_WRITE_PERMISSIONS); err == nil {
				util.OnExitError(file.Close)
				if self.Buffered {
					writer := util.NewBufferedWriter(file, BUFFER_SIZE)
					util.OnExitError(writer.Close)
					self.writer = writer
					klog.SetOutput(writer)
				} else {
					klog.SetOutput(util.NewSyncedWriter(file))
				}
			} else {
				util.Failf("log file error: %s", err.Error())
			}
		} else {
			if self.Buffered {
				writer := util.NewBufferedWriter(os.Stderr, BUFFER_SIZE)
				util.OnExitError(writer.Close)
				self.writer = writer
				klog.SetOutput(writer)
			} else {
				klog.SetOutput(util.NewSyncedWriter(os.Stderr))
			}
		}

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
		depth += 2
		return commonlog.NewUnstructuredMessage(func(message string) {
			if name := strings.Join(name, "."); len(name) > 0 {
				message = name + ": " + message
			}

			switch level {
			case commonlog.Critical:
				klog.ErrorDepth(depth, message)
			case commonlog.Error:
				klog.ErrorDepth(depth, message)
			case commonlog.Warning:
				klog.WarningDepth(depth, message)
			case commonlog.Notice:
				klog.InfoDepth(depth, message)
			case commonlog.Info:
				klog.InfoDepth(depth, message)
			case commonlog.Debug:
				klog.InfoDepth(depth, message)
			default:
				panic(fmt.Sprintf("unsupported level: %d", level))
			}
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
