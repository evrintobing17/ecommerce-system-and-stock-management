package shared

import (
	"log"
	"os"
	"time"
)

type Logger struct {
	*log.Logger
}

func NewLogger(prefix string) *Logger {
	return &Logger{
		log.New(os.Stdout, prefix, log.LstdFlags|log.Lmicroseconds|log.Lshortfile),
	}
}

func (l *Logger) RequestLog(method, path string, status int, duration time.Duration) {
	l.Printf("REQUEST: %s %s %d %v", method, path, status, duration)
}

func (l *Logger) ErrorLog(err error, context string) {
	l.Printf("ERROR: %s - %v", context, err)
}

func (l *Logger) InfoLog(message string) {
	l.Printf("INFO: %s", message)
}
