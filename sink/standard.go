package sink

import (
	"log"

	"github.com/tliron/commonlog"
)

func NewStandardLogger(parse LineParseFunc) *log.Logger {
	return log.New(NewPipeWriter(StandardLogParser), "", 0)
}

// TODO.
//
// Should take into account the logger flags. Example:
//
//	INFO 2023/10/21 11:15:46 simple_logger.go:73: Closing the StdScheduler.
//
// [LogParseFunc] signature
func StandardLogParser(line string) commonlog.Message {
	return nil
}
