package log

import (
	"log"
)

var logger Logger

func init() {
	logger = DefaultLog{}
}

func SetLogger(l Logger) {
	logger = l
}

func Debug(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}

func Info(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

func Warn(format string, args ...interface{}) {
	logger.Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}

type Logger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

type DefaultLog struct{}

func (s DefaultLog) Debugf(format string, args ...interface{}) {
	log.Printf(format+" \n", args...)
}

func (s DefaultLog) Infof(format string, args ...interface{}) {
	log.Printf(format+" \n", args...)
}

func (s DefaultLog) Warnf(format string, args ...interface{}) {
	log.Printf(format+" \n", args...)
}

func (s DefaultLog) Errorf(format string, args ...interface{}) {
	log.Printf(format+" \n", args...)
}
