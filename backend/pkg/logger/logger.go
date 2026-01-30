package logger

import (
	"log"
	"os"
)

type logger struct {
	info  *log.Logger
	error *log.Logger
	debug *log.Logger
}

var Logger *logger

func Init(level string) {
	Logger = &logger{
		info:  log.New(os.Stdout, "[INFO] ", log.LstdFlags),
		error: log.New(os.Stderr, "[ERROR] ", log.LstdFlags),
		debug: log.New(os.Stdout, "[DEBUG] ", log.LstdFlags),
	}
}

func Info(msg string) {
	if Logger != nil {
		Logger.info.Println(msg)
	}
}

func Error(msg string) {
	if Logger != nil {
		Logger.error.Println(msg)
	}
}

func Debug(msg string) {
	if Logger != nil {
		Logger.debug.Println(msg)
	}
}
