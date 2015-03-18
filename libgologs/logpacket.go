package libgologs

import (
	"time"
)

//======================================== log packet

type LogPacket struct {
	msgstr string
	msg    []interface{}
	time   time.Time
	level  int
}

func MakeLogPacket(level int, msg []interface{}) *LogPacket {
	return &LogPacket{level: level, msg: msg, time: time.Now().UTC()}
}
