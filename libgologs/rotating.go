package libgologs

import (
	"github.com/rshmelev/lumberjack"
	"io"
	"log"
	"strings"
	"time"
)

func CreateDailyRotatingLogger(Filename string, MaxSize int) *log.Logger {
	lu := CreateDailyRotatingWriteCloser(Filename, MaxSize)
	stdlog := log.New(lu, "", 0)
	return stdlog
}

func CreateDailyRotatingWriteCloser(Filename string, MaxSize int) io.WriteCloser {
	lu := SetupRotating(&lumberjack.Logger{
		Filename: Filename,
		MaxSize:  MaxSize,
	})
	return lu
}

func SetupRotating(lj *lumberjack.Logger) *lumberjack.Logger {

	// ensure that we'll rotate our log file
	go func() {
		t := time.Now().UTC()
		for t2 := range time.Tick(time.Second) {
			if t.Day() != t2.UTC().Day() {
				t = t2.UTC()
				log.Println("..." + t.Format("2006-01-02"))
				lj.Rotate()
			}
		}
	}()

	datetime := lumberjack.GetFirstDateTimeFromFile(lj.Filename)
	if datetime == "" {
		//println("log seems to not contain date info")
		return lj
	}

	currentdate := time.Now().UTC().Format("2006-01-02")

	date, _ := RegexExtract(".*(\\d{4}\\D\\d{2}\\D\\d{2}).*", datetime)
	date = strings.Replace(date, "/", "-", -1)               // normalize a bit
	currentdate = strings.Replace(currentdate, "/", "-", -1) // normalize a bit

	if currentdate == date {
		println("appending to current log file...")
	} else {
		println("rotating log: " + currentdate + " != " + date)
		lj.Rotate()
	}

	return lj
}
