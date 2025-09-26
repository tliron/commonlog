package sink

import (
	"bufio"
	"io"

	"github.com/tliron/commonlog"
	"github.com/tliron/go-kutil/util"
)

type LineParseFunc func(line string) commonlog.Message

func NewPipeWriter(parse LineParseFunc, name ...string) io.Writer {
	pipeReader, pipeWriter := io.Pipe()
	util.OnExitError(pipeReader.Close)

	go func() {
		scanner := bufio.NewScanner(pipeReader)
		for scanner.Scan() {
			if m := parse(scanner.Text()); m != nil {
				m.Send()
			}
		}

		if err := scanner.Err(); err != nil {
			if m := commonlog.NewCriticalMessage(0, name...); m != nil {
				m.Set(commonlog.MESSAGE, err.Error())
				m.Send()
			}
		}
	}()

	return pipeWriter
}
