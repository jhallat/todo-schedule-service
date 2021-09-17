package logger

import (
	"log"
	"os"
)

type logLevel int

const (
	DEBUG logLevel = 0 + iota
	INFO
	WARNING
	ERROR
	FATAL
)

var loggingLevel logLevel = DEBUG

func SetLoggLevel(level logLevel) {
	loggingLevel = level
}


func LogMessage(level logLevel, message string) {
	file, err := os.OpenFile("log.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)
	switch level {
	case INFO:
		logger := log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
		logger.Println(message)
	case WARNING:
		logger := log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
		logger.Println(message)
	case ERROR:
		logger := log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
		logger.Println(message)
	case FATAL:
		logger := log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
		logger.Fatal(message)
	}
}

func Debug(message string) {
	file, err := os.OpenFile("log.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)
	logger := log.New(file, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	logger.Println(message)
}

func Warn(message string) {
	file, err := os.OpenFile("log.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)
	logger := log.New(file, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile)
	logger.Println(message)
}


