package zerolog

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/rs/zerolog"
	logpkg "github.com/rs/zerolog/log"
	"github.com/tliron/commonlog"
	"github.com/tliron/kutil/terminal"
	"github.com/tliron/kutil/util"
)

const (
	LOG_FILE_WRITE_PERMISSIONS = 0600
	DEFAULT_BUFFER_SIZE        = 1_000
	TIME_FORMAT                = "2006/01/02 15:04:05.000"
)

func init() {
	backend := NewBackend()
	backend.Configure(0, nil)
	commonlog.SetBackend(backend)
}

//
// Backend
//

// Note: using commonlog to wrap zerolog will circumvent its primary optimization, which is for
// high performance and low resource use due to aggressively avoiding allocations. If you need
// that optimization then you should use zerolog's API directly.

type Backend struct {
	Writer     io.Writer
	BufferSize int
	Buffered   bool

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
		self.nameHierarchy.SetMaxLevel(commonlog.None)
		logpkg.Logger = zerolog.New(self.Writer)
		zerolog.SetGlobalLevel(zerolog.Disabled)
	} else {
		if path != nil {
			if file, err := os.OpenFile(*path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, LOG_FILE_WRITE_PERMISSIONS); err == nil {
				util.OnExitError(file.Close)
				if self.Buffered {
					writer := util.NewBufferedWriter(file, self.BufferSize, false)
					util.OnExitError(writer.Close)
					self.Writer = writer
				} else {
					self.Writer = util.NewSyncedWriter(file)
				}
				logpkg.Logger = zerolog.New(self.Writer)
			} else {
				util.Failf("log file error: %s", err.Error())
			}
		} else {
			self.Writer = os.Stderr
			if terminal.ColorizeStderr {
				// Note: ConsoleWriter has its own built-in support for
				// colorization, including for Windows terminals, which
				// relies on Out being equal to Stdout or Stderr, thus
				// we shouldn't use any wrappers for Out such as
				// BufferedWriter or SyncedWriter
				logpkg.Logger = zerolog.New(zerolog.ConsoleWriter{
					Out:        self.Writer,
					TimeFormat: TIME_FORMAT,
				})
			} else {
				logpkg.Logger = zerolog.New(zerolog.ConsoleWriter{
					Out:        self.Writer,
					TimeFormat: TIME_FORMAT,
					NoColor:    true,
				})
			}
		}

		zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMicro
		logpkg.Logger = logpkg.With().Timestamp().Logger()

		self.nameHierarchy.SetMaxLevel(maxLevel)
	}
}

// ([commonlog.Backend] interface)
func (self *Backend) GetWriter() io.Writer {
	return self.Writer
}

// ([commonlog.Backend] interface)
func (self *Backend) NewMessage(level commonlog.Level, depth int, name ...string) commonlog.Message {
	if self.AllowLevel(level, name...) {
		context := logpkg.With()
		if name := strings.Join(name, "."); len(name) > 0 {
			context = context.Str("name", name)
		}
		logger := context.Logger()

		var event *zerolog.Event
		switch level {
		case commonlog.Critical:
			event = logger.Error()
		case commonlog.Error:
			event = logger.Error()
		case commonlog.Warning:
			event = logger.Warn()
		case commonlog.Notice:
			event = logger.Info()
		case commonlog.Info:
			event = logger.Debug()
		case commonlog.Debug:
			event = logger.Trace()
		default:
			panic(fmt.Sprintf("unsupported log level: %d", level))
		}

		return commonlog.TraceMessage(NewMessage(event), depth)
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
