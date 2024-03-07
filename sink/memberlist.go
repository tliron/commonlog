package sink

import (
	"log"
	"strings"

	"github.com/hashicorp/memberlist"
	"github.com/tliron/commonlog"
)

const (
	MemberlistErrPrefix   = "[ERR] memberlist: "
	MemberlistWarnPrefix  = "[WARN] memberlist: "
	MemberlistDebugPrefix = "[DEBUG] memberlist: "
)

func NewMemberlistStandardLog(name ...string) *log.Logger {
	return NewStandardLogger(func(line string) commonlog.Message {
		level := commonlog.Debug

		if strings.HasPrefix(line, MemberlistErrPrefix) {
			line = line[len(MemberlistErrPrefix):]
			level = commonlog.Error
		} else if strings.HasPrefix(line, MemberlistWarnPrefix) {
			line = line[len(MemberlistWarnPrefix):]
			level = commonlog.Warning
		} else if strings.HasPrefix(line, MemberlistDebugPrefix) {
			line = line[len(MemberlistDebugPrefix):]
		}

		if message := commonlog.NewMessage(level, 2, name...); message != nil {
			message.Set(commonlog.MESSAGE, line)
			return message
		} else {
			return nil
		}
	})
}

//
// MemberlistEventLog
//

type MemberlistEventLog struct {
	log commonlog.Logger
}

func NewMemberlistEventLog(log commonlog.Logger) *MemberlistEventLog {
	return &MemberlistEventLog{log}
}

// ([memberlist.EventDelegate] interface)
func (self *MemberlistEventLog) NotifyJoin(node *memberlist.Node) {
	self.log.Infof("node has joined: %s", node.String())
}

// ([memberlist.EventDelegate] interface)
func (self *MemberlistEventLog) NotifyLeave(node *memberlist.Node) {
	self.log.Infof("node has left: %s", node.String())
}

// ([memberlist.EventDelegate] interface)
func (self *MemberlistEventLog) NotifyUpdate(node *memberlist.Node) {
	self.log.Infof("node was updated: %s", node.String())
}
