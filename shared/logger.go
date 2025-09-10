package shared

import (
	"log"
	"os"
	"runtime"
	"time"
)

type (
	ListErrors struct {
		Error    string
		File     string
		Function string
		Line     int
		Extra    interface{} `json:"extra,omitempty"`
	}
	Fields map[string]interface{}
)

type Log interface {
	SetMessageLog(err error, depthList ...int) *ListErrors
	RequestLog(method, path string, status int, duration time.Duration)
	ErrorLog(err error)
	InfoLog(message string)
}

type Logger struct {
	*log.Logger
}

func NewLogger(prefix string) Log {
	return &Logger{
		log.New(os.Stdout, prefix+" ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile),
	}
}

func (l *Logger) SetMessageLog(err error, depthList ...int) *ListErrors {
	var depth int
	if depthList == nil {
		depth = 1
	} else {
		depth = depthList[0]
	}
	le := new(ListErrors)
	if function, file, line, ok := runtime.Caller(depth); ok {
		le.Error = err.Error()
		le.File = file
		le.Function = runtime.FuncForPC(function).Name()
		le.Line = line
	} else {
		le = nil
	}
	return le
}

func (l *Logger) RequestLog(method, path string, status int, duration time.Duration) {
	l.Printf("REQUEST: %s %s %d %v", method, path, status, duration)
}

func (l *Logger) ErrorLog(err error) {
	l.Printf("ERROR: %v", err)
}

func (l *Logger) InfoLog(message string) {
	l.Printf("INFO: %s", message)
}
