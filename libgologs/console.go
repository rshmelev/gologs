package libgologs

import (
	"bytes"
)

func LogToConsole(p *LogPacket) {
	LogToConsole2(p, SliceToString(p.msg))
}
func LogToConsole2(p *LogPacket, msg string) {
	var buffer bytes.Buffer
	buffer.WriteString(p.time.Format("15:04:05.000 "))
	buffer.WriteString(LogLevels[p.level])
	buffer.WriteString(" ")
	buffer.WriteString(msg)
	buffer.WriteString("\n")

	print(buffer.String())
}

func NewConsoleLogger() SomeLogger {
	return NewLogger(LogDispatcherFunc(LogToConsole))
}
