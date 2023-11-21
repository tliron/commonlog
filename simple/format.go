package simple

import (
	"strings"
	"time"

	"github.com/tliron/commonlog"
	"github.com/tliron/kutil/terminal"
)

const TIME_FORMAT = "2006/01/02 15:04:05.000"

type FormatFunc func(message *commonlog.UnstructuredMessage, name []string, level commonlog.Level, colorize bool) string

// ([FormatFunc] signature)
func DefaultFormat(message *commonlog.UnstructuredMessage, name []string, level commonlog.Level, colorize bool) string {
	var s strings.Builder

	FormatLevel(&s, level, true)

	s.WriteRune(' ')
	FormatName(&s, name)

	if scope := message.ScopeString(); scope != "" {
		s.WriteRune(' ')
		s.WriteString(message.ScopeString())
	}

	s.WriteRune(' ')
	s.WriteString(message.Message)

	s_ := s.String()
	if colorize {
		s_ = FormatColorize(s_, level)
	}

	s = strings.Builder{}
	FormatTime(&s)
	s.WriteString(s_)

	if values := message.ValuesString(false); values != "" {
		if s.Len() > 0 {
			s.WriteRune(' ')
		}
		s.WriteString(values)
	}

	if location := message.LocationString(); location != "" {
		/*if colorize {
			location = FormatColorize(location, level)
		}*/
		s.WriteString("\n└─")
		s.WriteString(location)
	}

	return s.String()
}

func FormatTime(builder *strings.Builder) {
	builder.WriteString(time.Now().Format(TIME_FORMAT))
}

func FormatName(builder *strings.Builder, name []string) {
	builder.WriteRune('[')
	switch length := len(name); length {
	case 0:
	case 1:
		builder.WriteString(name[0])
	default:
		last := length - 1
		for _, n := range name[:last] {
			builder.WriteString(n)
			builder.WriteRune('.')
		}
		builder.WriteString(name[last])
	}
	builder.WriteRune(']')
}

func FormatLevel(writer *strings.Builder, level commonlog.Level, align bool) {
	if align {
		switch level {
		case commonlog.Critical:
			writer.WriteString("  CRIT")
		case commonlog.Error:
			writer.WriteString(" ERROR")
		case commonlog.Warning:
			writer.WriteString("  WARN")
		case commonlog.Notice:
			writer.WriteString("  NOTE")
		case commonlog.Info:
			writer.WriteString("  INFO")
		case commonlog.Debug:
			writer.WriteString(" DEBUG")
		}
	} else {
		switch level {
		case commonlog.Critical:
			writer.WriteString("CRIT")
		case commonlog.Error:
			writer.WriteString("ERROR")
		case commonlog.Warning:
			writer.WriteString("WARN")
		case commonlog.Notice:
			writer.WriteString("NOTE")
		case commonlog.Info:
			writer.WriteString("INFO")
		case commonlog.Debug:
			writer.WriteString("DEBUG")
		}
	}
}

func FormatColorize(s string, level commonlog.Level) string {
	switch level {
	case commonlog.Critical:
		return terminal.ColorRed(s)
	case commonlog.Error:
		return terminal.ColorRed(s)
	case commonlog.Warning:
		return terminal.ColorYellow(s)
	case commonlog.Notice:
		return terminal.ColorMagenta(s)
	case commonlog.Info:
		return terminal.ColorBlue(s)
	case commonlog.Debug:
		return terminal.ColorCyan(s)
	default:
		return s
	}
}
