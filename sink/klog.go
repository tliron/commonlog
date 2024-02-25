package sink

import (
	"strconv"
	"strings"

	"github.com/tliron/commonlog"
	"k8s.io/klog/v2"
)

// Example:
//
//	I0225 00:22:21.901297       1 handler.go:275] Adding GroupVersion tko.nephio.org v1alpha1 to ResourceManager

func CaptureKlogOutput(name ...string) {
	klog.LogToStderr(false)

	klog.SetOutput(NewPipeWriter(func(line string) commonlog.Message {
		severity := line[0]
		//month := line[1:3]
		//day := line[3:5]
		//hour := line[6:8]
		//minute := line[9:11]
		//second := line[12:14]
		//secondFraction := line[15:20]
		thread, _ := strconv.Atoi(line[21:29])
		line = line[30:]
		before, after, _ := strings.Cut(line, "] ")
		message := after
		before, after, _ = strings.Cut(before, ":")
		file := before
		lineNo, _ := strconv.Atoi(after)

		var level commonlog.Level
		switch severity {
		case 'I':
			level = commonlog.Info
		case 'W':
			level = commonlog.Warning
		case 'E':
			level = commonlog.Error
		case 'F':
			level = commonlog.Critical
		}

		if m := commonlog.NewMessage(level, 1, name...); m != nil {
			m.Set(commonlog.MESSAGE, message)

			if commonlog.Trace {
				m.Set(commonlog.FILE, file)
				m.Set(commonlog.LINE, lineNo)
				m.Set("_thread", thread)
			}

			return m
		} else {
			return nil
		}
	}))
}
