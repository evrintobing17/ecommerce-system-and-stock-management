package shared

import (
	"log"
	"os"
)

// InitLogger initializes the application logger
func InitLogger() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

// Info logs informational messages
func Info(format string, v ...interface{}) {
	log.Printf("INFO: "+format, v...)
}

// Error logs error messages
func Error(format string, v ...interface{}) {
	log.Printf("ERROR: "+format, v...)
}

// Debug logs debug messages
func Debug(format string, v ...interface{}) {
	log.Printf("DEBUG: "+format, v...)
}
