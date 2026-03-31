package utils

import (
	"log"
	"os"
)

type Logger struct {
	infoLogger 		*log.Logger
	errorLogger		*log.Logger
	warnLogger		*log.Logger
}

var logger *Logger


func InitLogger() {
	logger = &Logger{
		infoLogger: log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		errorLogger: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		warnLogger: log.New(os.Stdout, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func GetLogger() *Logger {
	if logger == nil {
		InitLogger()
	}
	return logger
}

func (l *Logger) Info(message string) {
	l.infoLogger.Println(message)
}

func (l *Logger) Error(message string) {
	l.errorLogger.Println(message)
}

func (l *Logger) Warn(message string) {
	l.warnLogger.Println(message)
}

func (l *Logger) Fatal(message string) {
	l.errorLogger.Fatal(message)
}