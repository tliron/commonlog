package sink

import (
	"bufio"
	"io"
	"log"

	"github.com/tliron/commonlog"
	"github.com/tliron/kutil/util"
)

type StandardLogParseFunc func(line string) commonlog.Message

func NewStandardLogger(parse StandardLogParseFunc) *log.Logger {
	pipeReader, pipeWriter := io.Pipe()
	util.OnExitError(pipeReader.Close)

	go func() {
		reader := bufio.NewReader(pipeReader)
		for {
			if line, err := reader.ReadString('\n'); err == nil {
				if len(line) > 1 {
					line = line[:len(line)-1]
					if message := parse(line); message != nil {
						message.Send()
					}
				}
			} else {
				return
			}
		}
	}()

	return log.New(pipeWriter, "", 0)
}

// TODO.
//
// Should take into account the logger flags. Example:
//
//	INFO 2023/10/21 11:15:46 simple_logger.go:73: Closing the StdScheduler.
//
// [StandardLogParseFunc] signature
func DefaultStandardLogParser(line string) commonlog.Message {
	return nil
}
