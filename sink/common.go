package sink

import (
	"bufio"
	"io"

	"github.com/tliron/commonlog"
	"github.com/tliron/kutil/util"
)

type LineParseFunc func(line string) commonlog.Message

func NewPipeWriter(parse LineParseFunc, name ...string) io.Writer {
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

	return pipeWriter
}
