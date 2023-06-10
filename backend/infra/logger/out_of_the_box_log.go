package logger

var _logger Logger

func GetLogger() Logger {
	return _logger
}

func SetLogger(l Logger) {
	_logger = l
}

func Debug(args ...interface{}) {
	_logger.AddCallerSkip(1).Debug(args...)
}

func Info(args ...interface{}) {
	_logger.AddCallerSkip(1).Info(args...)
}

func Warn(args ...interface{}) {
	_logger.AddCallerSkip(1).Warn(args...)
}

func Error(args ...interface{}) {
	_logger.AddCallerSkip(1).Error(args...)
}

func Panic(args ...interface{}) {
	_logger.AddCallerSkip(1).Panic(args...)
}

func Fatal(args ...interface{}) {
	_logger.AddCallerSkip(1).Fatal(args...)
}

func Debugf(format string, args ...interface{}) {
	_logger.AddCallerSkip(1).Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	_logger.AddCallerSkip(1).Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	_logger.AddCallerSkip(1).Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	_logger.AddCallerSkip(1).Errorf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	_logger.AddCallerSkip(1).Panicf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	_logger.AddCallerSkip(1).Fatalf(format, args...)
}

func Named(name string) Logger {
	return _logger.Named(name)
}

func With(fields ...Field) Logger {
	return _logger.With(fields...)
}

func GetLevel() string {
	return _logger.GetLevel()
}

func SetLevel(level string) {
	_logger.SetLevel(level)
}

func Sync() error {
	return _logger.Sync()
}

func init() {
	_logger = GetLoggerFactory().CreateDefaultLogger()
}
