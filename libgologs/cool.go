package libgologs

import (
	"log"
	"sync/atomic"
	"time"
)

//=======================================

type WebsocketLogMsg struct {
	L int
	M string
	T int64
}

type WebsocketLogFunc func(*WebsocketLogMsg)

type CoolLogger struct {
	LoggerBase
	MemoryLimit         int
	WebsocketFunc       WebsocketLogFunc
	FullLogFilename     string
	dailyRotatingLogger *log.Logger

	historyChannel chan []*WebsocketLogMsg
}

func (p *CoolLogger) GetHistory() []*WebsocketLogMsg {
	p.DispatcherFunc(MakeLogPacket(LLEVEL_HISTORY, nil))
	history := <-p.historyChannel
	return history
}

func NewCoolLogger(p *CoolLogger) SomeLogger {

	// i like an idea to retrieve logs via websockets in realtime.
	// it simplifies debugging on remote server for me.
	// i know there are other ways to get this functionality

	p.historyChannel = make(chan []*WebsocketLogMsg)
	history := make([]*WebsocketLogMsg, 1)
	history[0] = &WebsocketLogMsg{L: LLEVEL_STARTED, M: "", T: time.Now().UTC().UnixNano() / int64(time.Millisecond)}

	logch := make(chan *LogPacket, 1024)
	flushch := make(chan int, 10)

	if p.FullLogFilename != "" {
		p.dailyRotatingLogger = CreateDailyRotatingLogger(p.FullLogFilename, 500) // 500 MB
	}
	var lastSecond int32 = 0

	p.LoggerBase.DispatcherFunc = func(pack *LogPacket) {
		if pack.level >= 0 {
			pack.msgstr = SliceToString(pack.msg)

			currentSecond := int32(pack.time.Second())
			aLastSecond := atomic.LoadInt32(&lastSecond)
			var lastSecondChanged bool
			if lastSecond != currentSecond {
				lastSecondChanged = atomic.CompareAndSwapInt32(&lastSecond, aLastSecond, currentSecond)
				if lastSecondChanged {
					// i like the idea to have newline before each new second
					// that helps a lot in understanding app events timings
					println("")
				}
			}

			// console logging should happend in main thread
			LogToConsole2(pack, pack.msgstr)

			if lastSecondChanged && p.dailyRotatingLogger != nil {
				p.dailyRotatingLogger.Println("")
			}

		}
		logch <- pack
		if pack.level == LLEVEL_FLUSH { // wait for flush finish
			<-flushch
		}
	}

	go func() {
		for {
			pack := <-logch
			if pack.level == LLEVEL_FLUSH {
				// flush if needed ...
				flushch <- 1
				continue
			}
			if pack.level == LLEVEL_HISTORY {
				p.historyChannel <- history
				continue
			}

			wsmsg := &WebsocketLogMsg{L: pack.level, M: pack.msgstr, T: pack.time.UnixNano() / int64(time.Millisecond)}

			if p.MemoryLimit > 0 {
				history = append(history, wsmsg)
				over := len(history) - p.MemoryLimit
				if over > 0 {
					history = history[over:]
				}
			}

			if p.WebsocketFunc != nil {
				go p.WebsocketFunc(wsmsg)
			}

			if p.dailyRotatingLogger != nil {
				p.dailyRotatingLogger.Println(pack.time.Format("2006-01-02 15:04:05.000"), LogLevelsShort[pack.level], pack.msgstr)
			}
		}
	}()

	return p
}
