package czlog

import (
	"log"
	"gopkg.in/natefinch/lumberjack.v2"
	"fmt"
	"github.com/ivpusic/grpool"
)

const (
	LevelError = iota
	LevelWarn
	LevelInfo
	LevelDebug
)

const (
	Ldate         = 1 << iota     // the date in the local time zone: 2009/01/23
	Ltime                         // the time in the local time zone: 01:23:23
	Lmicroseconds                 // microsecond resolution: 01:23:23.123123.  assumes Ltime.
	Llongfile                     // full file name and line number: /a/b/c/d.go:23
	Lshortfile                    // final file name element and line number: d.go:23. overrides Llongfile
	LUTC                          // if Ldate or Ltime is set, use UTC rather than the local time zone
	LstdFlags     = Ldate | Ltime // initial values for the standard logger
)

type Logger struct {
	level int
	err     *log.Logger
	warn	*log.Logger
	info 	*log.Logger
	debug 	*log.Logger
	p 		*grpool.Pool
	depth	int
	compress 	bool
	localTime	bool
}


func New(filename string, maxSize int, flag int, compress bool, localTime bool) *Logger {
	logger := new(Logger)


	jack := &lumberjack.Logger{
		Filename: filename,
		MaxSize: maxSize,
		LocalTime: localTime,
		Compress: compress,
	}

	logger.err = log.New(jack, "[ERROR] ", flag)
	logger.warn = log.New(jack, "[WARN] ", flag)
	logger.info = log.New(jack, "[INFO] ", flag)
	logger.debug = log.New(jack, "[DEBUG] ", flag)
	logger.p = grpool.NewPool(1,50)

	logger.depth = 3

	return logger
}

//var std = New("czloggertst.log", 1,log.LstdFlags)

func (ll *Logger) SetLevel(l int) {
	ll.level = l
}

func (ll *Logger) Error(format string, v ...interface{}) {
	//fmt.Printf("testing0...\n")
	if LevelError > ll.level {
		return
	}
	//fmt.Printf("testing1...\n")
	ll.p.JobQueue <- func() {
		//fmt.Printf("tesing2.1...\n")
		ll.err.Output(ll.depth, fmt.Sprintf(format, v...))
		//fmt.Printf("testing2...\n")
	}
	//fmt.Printf("Jobquene cap: %d, length: %d\n", cap(ll.p.JobQueue), len(ll.p.JobQueue))
}

func (ll *Logger) Warn(format string, v ...interface{}) {
	if LevelWarn > ll.level {
		return
	}
	ll.p.JobQueue <- func() {
		ll.warn.Output(ll.depth, fmt.Sprintf(format, v...))
	}
}

func (ll *Logger) Info(format string, v ...interface{}) {
	if LevelInfo > ll.level {
		return
	}
	ll.p.JobQueue <- func() {
		ll.info.Output(ll.depth, fmt.Sprintf(format, v...))
	}
}

func (ll *Logger) Debug(format string, v ...interface{}) {
	if LevelDebug > ll.level {
		return
	}
	ll.p.JobQueue <- func() {
		ll.debug.Output(ll.depth, fmt.Sprintf(format, v...))
	}
}

func (ll *Logger) ReleasePool() {
	ll.p.Release()
}