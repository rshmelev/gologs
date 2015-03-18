package libgologs

import (
	"log"
)

type SomeLogger interface {
	Debug(msg ...interface{})
	Info(msg ...interface{})
	Warn(msg ...interface{})
	Error(msg ...interface{})
	Dispatch(level int, msg ...interface{})
	Flush()
	SetAsStdLogWriter(flags ...int)
}

type LogDispatcherFunc func(*LogPacket)
type LoggerBase struct {
	DispatcherFunc LogDispatcherFunc
}

var LogLevels = []string{"TRACE", "DEBUG", "INFO", "WARN", "ERROR", "FATAL"}
var LogLevelsShort = []string{"T", "D", "I", "W", "E", "F"}

const (
	LLEVEL_UNKNOWN = iota
	LLEVEL_DEBUG   = iota
	LLEVEL_INFO    = iota
	LLEVEL_WARN    = iota
	LLEVEL_ERROR   = iota

	LLEVEL_FLUSH   = -1
	LLEVEL_HISTORY = -2
	LLEVEL_STARTED = -3
)

//==================================================================

func NewLogger(disp LogDispatcherFunc) SomeLogger {
	return SomeLogger(&LoggerBase{disp})
}

//------

func (base *LoggerBase) Debug(msg ...interface{}) {
	base.DispatcherFunc(MakeLogPacket(LLEVEL_DEBUG, msg))
}
func (base *LoggerBase) Info(msg ...interface{}) {
	base.DispatcherFunc(MakeLogPacket(LLEVEL_INFO, msg))
}
func (base *LoggerBase) Warn(msg ...interface{}) {
	base.DispatcherFunc(MakeLogPacket(LLEVEL_WARN, msg))
}
func (base *LoggerBase) Error(msg ...interface{}) {
	base.DispatcherFunc(MakeLogPacket(LLEVEL_ERROR, msg))
}
func (base *LoggerBase) Dispatch(level int, msg ...interface{}) {
	base.DispatcherFunc(MakeLogPacket(level, msg))
}

func (base *LoggerBase) Flush() {
	base.DispatcherFunc(MakeLogPacket(LLEVEL_FLUSH, nil))
}

func (base *LoggerBase) SetAsStdLogWriter(flags ...int) {
	log.SetPrefix("")
	log.SetOutput(base)
	if len(flags) > 0 {
		log.SetFlags(flags[0])
	} else {
		log.SetFlags(log.Lshortfile)
	}
}

func (base *LoggerBase) Write(p []byte) (n int, err error) {
	lenp := len(p)

	w := p
	if lenp > 0 && p[lenp-1] == '\n' {
		w = p[:lenp-1]
	}
	line, level := DetectLevelAndCutPrefix(w, true)
	base.Dispatch(level, string(line))

	return lenp, nil
}

// fastest possible func to detect the level and cut level prefix
func DetectLevelAndCutPrefix(b []byte, cutprefix bool) ([]byte, int) {
	blen := len(b)
	if b == nil || len(b) < 4 {
		return b, LLEVEL_INFO
	}
	res := LLEVEL_INFO
	if b[0] == 'E' && b[1] == 'R' && b[2] == 'R' {
		res = LLEVEL_ERROR
	} else if b[0] == 'W' && b[1] == 'A' && b[2] == 'R' && b[3] == 'N' {
		res = LLEVEL_WARN
	} else if b[0] == 'D' && b[1] == 'E' && b[2] == 'B' && b[3] == 'U' {
		res = LLEVEL_DEBUG
	}
	if res != LLEVEL_INFO && cutprefix {
		state := 0
		for i := 3; i < blen; i++ {
			if state == 0 {
				if b[i] == ':' {
					state = 1
				}
			} else if b[i] != ' ' {
				return b[i:], res
			}
		}
	}
	return b, res
}
