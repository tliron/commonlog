package commonlog

import (
	"fmt"
	"io"
	"strings"

	"github.com/tliron/kutil/terminal"
)

var backend Backend

func SetBackend(backend_ Backend) {
	backend = backend_
}

func Configure(verbosity int, path *string) {
	if backend != nil {
		backend.Configure(verbosity, path)
	}
}

func Initialize(verbosity int, path string) {
	if path == "" {
		if terminal.Quiet {
			verbosity = -4
		}
		Configure(verbosity, nil)
	} else {
		Configure(verbosity, &path)
	}
}

func GetWriter() io.Writer {
	if backend != nil {
		return backend.GetWriter()
	} else {
		return nil
	}
}

func AllowLevel(name []string, level Level) bool {
	if backend != nil {
		return backend.AllowLevel(name, level)
	} else {
		return false
	}
}

func SetMaxLevel(name []string, level Level) {
	if backend != nil {
		backend.SetMaxLevel(name, level)
	}
}

func GetMaxLevel(name []string) Level {
	if backend != nil {
		return backend.GetMaxLevel(name)
	} else {
		return None
	}
}

func NewMessage(name []string, level Level, depth int) Message {
	if backend != nil {
		return backend.NewMessage(name, level, depth)
	} else {
		return nil
	}
}

func GetLogger(name string) Logger {
	name_ := strings.Split(name, ".")
	if len(name_) == 0 {
		name_ = nil
	}
	return NewBackendLogger(name_)
}

func GetLoggerf(format string, values ...any) Logger {
	return GetLogger(fmt.Sprintf(format, values...))
}

func CallAndLogError(f func() error, task string, log Logger) {
	if err := f(); err != nil {
		log.Errorf("%s: %s", task, err.Error())
	}
}
