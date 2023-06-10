package zapwrap

import (
	"io"
	"os"

	commonlog "github.com/labstack/gommon/log"

	log "code-shooting/infra/logger"
)

type wrappedLogger struct {
	zap log.Logger
}

func (l *wrappedLogger) SetHeader(h string) {
}

func (l *wrappedLogger) Output() io.Writer {
	return os.Stdout
}

func (l *wrappedLogger) SetOutput(w io.Writer) {

}
func (l *wrappedLogger) Prefix() string {
	return "echowrap"
}
func (l *wrappedLogger) SetPrefix(p string) {
}
func (l *wrappedLogger) Level() commonlog.Lvl {
	return commonlog.INFO
}
func (l *wrappedLogger) SetLevel(v commonlog.Lvl) {

}
func (l *wrappedLogger) Print(i ...interface{}) {
	l.zap.Info(i)
}
func (l *wrappedLogger) Printf(format string, args ...interface{}) {
	l.zap.Infof(format, args)
}
func (l *wrappedLogger) Printj(j commonlog.JSON) {
	l.Print(j)
}
func (l *wrappedLogger) Debug(i ...interface{}) {
	l.zap.Debug(i)
}
func (l *wrappedLogger) Debugf(format string, args ...interface{}) {
	l.zap.Debugf(format, args)
}
func (l *wrappedLogger) Debugj(j commonlog.JSON) {
	l.Debug(j)
}
func (l *wrappedLogger) Info(i ...interface{}) {
	l.Print(i)
}
func (l *wrappedLogger) Infof(format string, args ...interface{}) {
	l.Printf(format, args)
}
func (l *wrappedLogger) Infoj(j commonlog.JSON) {

}
func (l *wrappedLogger) Warn(i ...interface{}) {
	l.zap.Warn(i)
}
func (l *wrappedLogger) Warnf(format string, args ...interface{}) {
	l.zap.Warnf(format, args)
}
func (l *wrappedLogger) Warnj(j commonlog.JSON) {
	l.Warn(j)
}
func (l *wrappedLogger) Error(i ...interface{}) {
	l.zap.Error(i)
}
func (l *wrappedLogger) Errorf(format string, args ...interface{}) {
	l.zap.Errorf(format, args)
}
func (l *wrappedLogger) Errorj(j commonlog.JSON) {
	l.Error(j)
}
func (l *wrappedLogger) Fatal(i ...interface{}) {
	l.zap.Error(i)
}
func (l *wrappedLogger) Fatalj(j commonlog.JSON) {
	l.Error(j)
}
func (l *wrappedLogger) Fatalf(format string, args ...interface{}) {
	l.zap.Errorf(format, args)
}
func (l *wrappedLogger) Panic(i ...interface{}) {
	l.zap.Error(i)
}
func (l *wrappedLogger) Panicj(j commonlog.JSON) {
	l.Error(j)
}
func (l *wrappedLogger) Panicf(format string, args ...interface{}) {
	l.zap.Errorf(format, args)
}

func Wrap() *wrappedLogger {
	return &wrappedLogger{zap: log.GetLogger().AddCallerSkip(1)}
}
