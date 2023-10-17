package slog

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/tliron/commonlog"
	"github.com/tliron/kutil/util"
)

const LOG_FILE_WRITE_PERMISSIONS = 0600

const DEFAULT_BUFFER_SIZE = 1_000

func init() {
	backend := NewBackend()
	backend.Configure(0, nil)
	commonlog.SetBackend(backend)
}

//
// Backend
//

type Backend struct {
	Logger     *slog.Logger
	Writer     io.Writer
	BufferSize int
	Buffered   bool
	AddSource  bool

	nameHierarchy *commonlog.NameHierarchy
}

func NewBackend() *Backend {
	return &Backend{
		BufferSize:    DEFAULT_BUFFER_SIZE,
		Buffered:      true,
		nameHierarchy: commonlog.NewNameHierarchy(),
	}
}

// ([commonlog.Backend] interface)
func (self *Backend) Configure(verbosity int, path *string) {
	maxLevel := commonlog.VerbosityToMaxLevel(verbosity)

	if maxLevel == commonlog.None {
		self.Writer = io.Discard
		self.Logger = slog.New(MOCK_HANDLER)
		self.nameHierarchy.SetMaxLevel(commonlog.None)
	} else {
		if path != nil {
			if file, err := os.OpenFile(*path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, LOG_FILE_WRITE_PERMISSIONS); err == nil {
				util.OnExitError(file.Close)
				if self.Buffered {
					// Note: slog.NewTextHandler modifies its buffers, so we must copy byte slices
					writer := util.NewBufferedWriter(file, self.BufferSize, true)
					util.OnExitError(writer.Close)
					self.Writer = writer
				} else {
					self.Writer = util.NewSyncedWriter(file)
				}
			} else {
				util.Failf("log file error: %s", err.Error())
			}
		} else if self.Buffered {
			// Note: slog.NewTextHandler modifies its buffers, so we must copy byte slices
			writer := util.NewBufferedWriter(os.Stderr, self.BufferSize, true)
			util.OnExitError(writer.Close)
			self.Writer = writer
		} else {
			self.Writer = util.NewSyncedWriter(os.Stderr)
		}

		self.Logger = slog.New(slog.NewTextHandler(self.Writer, &slog.HandlerOptions{
			AddSource: self.AddSource,
			Level:     slog.LevelDebug,
		}))

		self.nameHierarchy.SetMaxLevel(maxLevel)
	}

	slog.SetDefault(self.Logger)
}

// ([commonlog.Backend] interface)
func (self *Backend) GetWriter() io.Writer {
	return self.Writer
}

// ([commonlog.Backend] interface)
func (self *Backend) NewMessage(level commonlog.Level, depth int, name ...string) commonlog.Message {
	if (self.Logger != nil) && self.AllowLevel(level, name...) {
		var slogLevel slog.Level
		switch level {
		case commonlog.Critical:
			slogLevel = slog.LevelError
		case commonlog.Error:
			slogLevel = slog.LevelError
		case commonlog.Warning:
			slogLevel = slog.LevelWarn
		case commonlog.Notice:
			slogLevel = slog.LevelInfo
		case commonlog.Info:
			slogLevel = slog.LevelInfo
		case commonlog.Debug:
			slogLevel = slog.LevelDebug
		default:
			panic(fmt.Sprintf("unsupported level: %d", level))
		}

		return NewMessage(self.Logger, slogLevel, context.Background())
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
