package logger

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/fatih/color"
)

var myLogger Logger

func init() {
	// Init the logger during package initialization
	log, err := NewLogger(0, 0, 0)
	if err != nil {
		panic(err)
	}

	myLogger = log

}

func GetLogger() Logger {
	return myLogger
}

type Logger interface {
	Info(format string, args ...interface{})
	Warning(format string, args ...interface{})
	Error(format string, args ...interface{})
	Debug(format string, args ...interface{})
}

// defaultLogger is the default implementation of the Logger interface.
type defaultLogger struct {
	infoLogger    *log.Logger
	warningLogger *log.Logger
	errorLogger   *log.Logger
	debugLogger   *log.Logger
	mu            sync.Mutex
}

// NewLogger creates a new instance of the defaultLogger with optional configurations.
func NewLogger(infoFlags, warningFlags, errorFlags int) (Logger, error) {
	infoLogger := log.New(os.Stdout, color.GreenString("[INF] "), infoFlags)
	if infoLogger == nil {
		return nil, fmt.Errorf("failed to create info logger")
	}

	warningLogger := log.New(os.Stdout, color.YellowString("[WARN] "), warningFlags)
	if warningLogger == nil {
		return nil, fmt.Errorf("failed to create warning logger")
	}

	errorLogger := log.New(os.Stdout, color.RedString("[ERROR] "), errorFlags)
	if errorLogger == nil {
		return nil, fmt.Errorf("failed to create error logger")
	}

	debugLogger := log.New(os.Stdout, color.CyanString("[DEBUG] "), log.Ldate|log.Ltime|log.Lshortfile)
	if debugLogger == nil {
		return nil, fmt.Errorf("failed to create debug logger")
	}

	return &defaultLogger{
		infoLogger:    infoLogger,
		warningLogger: warningLogger,
		errorLogger:   errorLogger,
		debugLogger:   debugLogger,
	}, nil
}

func (l *defaultLogger) Info(format string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	message := fmt.Sprintf(format, args...)
	l.infoLogger.Println(message)

}

func (l *defaultLogger) Warning(format string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	message := fmt.Sprintf(format, args...)
	l.warningLogger.Println(message)
}

func (l *defaultLogger) Error(format string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	message := fmt.Sprintf(format, args...)
	l.errorLogger.Println(message)
}

func (l *defaultLogger) Debug(format string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	message := fmt.Sprintf(format, args...)
	l.debugLogger.Println(message)
}

