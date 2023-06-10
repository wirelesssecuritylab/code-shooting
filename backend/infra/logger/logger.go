package logger

import (
	"log"
)

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Panic(args ...interface{})
	Fatal(args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Panicf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})

	CreateStdLogger() *log.Logger

	AddCallerSkip(skip int) Logger
	Named(name string) Logger
	With(fields ...Field) Logger
	SetLevel(level string)
	GetLevel() string
	Sync() error
	UpdateWriteSyncs() error
}
