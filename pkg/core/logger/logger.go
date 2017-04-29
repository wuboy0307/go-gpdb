package logger

import (
	"../../../src/logrus"
)

// Initalize the logger
func InitLogger() {
	var textFormatter = new(logrus.TextFormatter)
	textFormatter.TimestampFormat = "2006-01-02 15:04:05"
	textFormatter.FullTimestamp = true
	logrus.SetFormatter(textFormatter)
}


// Format the logger to print message
func Printf(format string, v ...interface{}) {
	logrus.Printf(format, v)
}

func Println(format string) {
	logrus.Println(format)
}

// Format the logger to FATAL message
func Fatalf(format string, v error) {
	logrus.Fatal(v)
}

func Fatal(format string) {
	logrus.Fatal(format)
}

// Format the logger to ERROR message
func Error(format string) {
	logrus.Error(format)
}

func Errorf(format string, v ...interface{}) {
	logrus.Error(format, v)
}


// Format the logger to WARNING message
func Warn(format string) {
	logrus.Warn(format)
}

func Warnf(format string, v ...interface{}) {
	logrus.Warn(format, v)
}