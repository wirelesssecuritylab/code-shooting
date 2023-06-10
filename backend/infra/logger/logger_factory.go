package logger

import (
	"fmt"
	"log"
	"sync"
	"time"

	"code-shooting/infra/common"
	"code-shooting/infra/config"
	"code-shooting/infra/config/model"
	"code-shooting/infra/logger/internal"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(configPath string, options ...Option) (Logger, error) {
	conf, err := config.NewConfig(configPath)
	if err != nil {
		return nil, errors.Wrap(err, "build config")
	}
	var namespace string
	if err = conf.Get(common.ROOT+"."+common.NODE+"."+common.APP_NAMESPACE, &namespace); err != nil && !config.IsNotExist(err) {
		return nil, errors.Wrap(err, "get namespace")
	}
	internal.SetNamespace(namespace)
	cfg, err := getConfigParsingService().Parse(conf, options...)
	if err != nil {
		return nil, errors.Wrap(err, "config content")
	}
	return GetLoggerFactory().CreateLogger(cfg, options...)
}

type loggerFactory interface {
	CreateLogger(cfg internal.Config, options ...Option) (Logger, error)
	CreateDefaultLogger() Logger
	createZapCore(cfg internal.Config, atomicLevel zap.AtomicLevel, options ...Option) (zapcore.Core, []func(), error)
}

func GetLoggerFactory() loggerFactory {
	once.Do(func() {
		loggerFactorySingleton = &loggerFactoryImpl{}
	})
	return loggerFactorySingleton
}

var once sync.Once
var loggerFactorySingleton loggerFactory

type loggerFactoryImpl struct {
}

func (s *loggerFactoryImpl) CreateLogger(cfg internal.Config, options ...Option) (Logger, error) {
	zapLevel, err := cfg.Level.ToZapLevel()
	if err != nil {
		return nil, errors.Wrap(err, "transform level")
	}
	atomicLevel := zap.NewAtomicLevelAt(zapLevel)
	core, closers, err := s.createZapCore(cfg, atomicLevel, options...)
	if err != nil {
		return nil, errors.Wrap(err, "create zapCore")
	}
	zapOptions := []zap.Option{zap.Development(), zap.AddCaller(), zap.AddCallerSkip(1)}
	Options(options).Do(&zapOptions)
	l := &logger{level: &atomicLevel, zLogger: zap.New(core, zapOptions...), cfg: cfg, opts: options, closers: closers}
	l.registerConfigChangeEvent()
	return l, nil
}

func (s *loggerFactoryImpl) createZapCore(cfg internal.Config, atomicLevel zap.AtomicLevel, options ...Option) (zapcore.Core, []func(), error) {
	zapEncoder, err := s.createZapEncoder(cfg.Encoder, cfg.Format, options...)
	if err != nil {
		return nil, nil, errors.Wrap(err, "create encoder")
	}
	core := zapcore.NewCore(zapEncoder, zapcore.NewMultiWriteSyncer(), atomicLevel)
	return core, nil, nil

}

func (s *loggerFactoryImpl) CreateDefaultLogger() Logger {
	atomicLevel := zap.NewAtomicLevelAt(zap.InfoLevel)
	core := zapcore.NewCore(s.createPlainEncoder(internal.CustomFormat{}), zapcore.NewMultiWriteSyncer(), atomicLevel)
	zapOptions := []zap.Option{zap.Development(), zap.AddCaller(), zap.AddCallerSkip(1)}
	return &logger{level: &atomicLevel, zLogger: zap.New(core, zapOptions...)}
}


func (s *loggerFactoryImpl) createZapEncoder(encoder internal.Encoder, format internal.CustomFormat, options ...Option) (zapcore.Encoder, error) {
	builders := map[internal.Encoder]func() zapcore.Encoder{
		internal.JsonEncoder: func() zapcore.Encoder {
			return s.createJsonEncoder()
		},
		internal.PlainEncoder: func() zapcore.Encoder {
			return s.createPlainEncoder(format, options...)
		},
	}

	for enc, builder := range builders {
		if enc.Equal(encoder) {
			return builder(), nil
		}
	}
	return nil, errors.Errorf("failed create encoder by %s", encoder.String())
}

func (s *loggerFactoryImpl) createPlainEncoder(format internal.CustomFormat, options ...Option) zapcore.Encoder {
	cfg := zapcore.EncoderConfig{
		TimeKey:        "T",
		EncodeTime:     internal.TimeEncoderOfLayout("2006-01-02T15:04:05.000Z07:00"),
		LevelKey:       "L",
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		NameKey:        "N",
		EncodeName:     zapcore.FullNameEncoder,
		CallerKey:      "C",
		EncodeCaller:   internal.PlainCallerEncoder,
		MessageKey:     "M",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeDuration: zapcore.StringDurationEncoder, // control the format of zap.Duration Field
	}
	Options(options).Do(&cfg)
	return internal.NewPlainEncoder(cfg, format)
}

func (s *loggerFactoryImpl) createJsonEncoder() zapcore.Encoder {
	cfg := zapcore.EncoderConfig{
		TimeKey:        "T",
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		LevelKey:       "L",
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		NameKey:        "N",
		EncodeName:     zapcore.FullNameEncoder,
		CallerKey:      "C",
		EncodeCaller:   zapcore.ShortCallerEncoder,
		MessageKey:     "M",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeDuration: zapcore.StringDurationEncoder, // control the format of zap.Duration Field
	}
	return zapcore.NewJSONEncoder(cfg)
}

type logger struct {
	level   *zap.AtomicLevel
	zLogger *zap.Logger
	zFields []zap.Field
	cfg     internal.Config
	opts    []Option
	closers []func()
}

func (s *logger) Debug(args ...interface{}) {
	s.zLogger.Debug(fmt.Sprint(args...), s.zFields...)
}

func (s *logger) Debugf(format string, args ...interface{}) {
	s.zLogger.Debug(fmt.Sprintf(format, args...), s.zFields...)
}

func (s *logger) Info(args ...interface{}) {
	s.zLogger.Info(fmt.Sprint(args...), s.zFields...)
}

func (s *logger) Infof(format string, args ...interface{}) {
	s.zLogger.Info(fmt.Sprintf(format, args...), s.zFields...)
}

func (s *logger) Warn(args ...interface{}) {
	s.zLogger.Warn(fmt.Sprint(args...), s.zFields...)
}

func (s *logger) Warnf(format string, args ...interface{}) {
	s.zLogger.Warn(fmt.Sprintf(format, args...), s.zFields...)
}

func (s *logger) Error(args ...interface{}) {
	s.zLogger.Error(fmt.Sprint(args...), s.zFields...)
}

func (s *logger) Errorf(format string, args ...interface{}) {
	s.zLogger.Error(fmt.Sprintf(format, args...), s.zFields...)
}

func (s *logger) Panic(args ...interface{}) {
	if s.level.Enabled(zapcore.PanicLevel) {
		s.zLogger.Panic(fmt.Sprint(args...), s.zFields...)
	}
}

func (s *logger) Panicf(format string, args ...interface{}) {
	if s.level.Enabled(zapcore.PanicLevel) {
		s.zLogger.Panic(fmt.Sprintf(format, args...), s.zFields...)
	}
}

func (s *logger) Fatal(args ...interface{}) {
	if s.level.Enabled(zapcore.FatalLevel) {
		s.zLogger.Fatal(fmt.Sprint(args...), s.zFields...)
	}
}

func (s *logger) Fatalf(format string, args ...interface{}) {
	if s.level.Enabled(zapcore.FatalLevel) {
		s.zLogger.Fatal(fmt.Sprintf(format, args...), s.zFields...)
	}
}

func (s *logger) CreateStdLogger() *log.Logger {
	return zap.NewStdLog(s.zLogger)
}

func (s *logger) AddCallerSkip(skip int) Logger {
	return &logger{level: s.level, zLogger: s.zLogger.WithOptions(zap.AddCallerSkip(skip))}
}

func (s *logger) Named(name string) Logger {
	return &logger{level: s.level, zLogger: s.zLogger.Named(name)}
}

func (s *logger) With(fields ...Field) Logger {
	l := &logger{level: s.level, zLogger: s.zLogger,
		zFields: make([]zapcore.Field, 0, len(s.zFields)+len(fields))}

	for _, field := range s.zFields {
		l.zFields = append(l.zFields, field)
	}

	for _, field := range fields {
		l.zFields = append(l.zFields, field())
	}

	return l
}

func (s *logger) SetLevel(level string) {
	s.cfg.Level = internal.Level(level)
	zapLevel, err := s.cfg.Level.ToZapLevel()
	if err != nil {
		s.Error("set log level err: ", err.Error())
		return
	}
	s.level.SetLevel(zapLevel)
}

func (s *logger) UpdateWriteSyncs() error {
	core, closers, err := GetLoggerFactory().createZapCore(s.cfg, *s.level, s.opts...)
	if err != nil {
		return errors.Wrap(err, " update write syncs")
	}
	c := s.zLogger.WithOptions(zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return core
	}))
	cs := s.closers
	time.AfterFunc(10*time.Second, func() {
		for _, close := range cs {
			close()
		}
	})

	s.zLogger = c
	s.closers = closers
	return nil
}

func (s *logger) GetLevel() string {
	return s.level.String()
}

func (s *logger) Sync() error {
	return s.zLogger.Sync()
}

func (s *logger) registerConfigChangeEvent() {
	s.registerLogLevelEventHandler()
	s.registerNamespaceEventHandler()
	s.registerRotateConfigEventHandler()
}

func (s *logger) registerLogLevelEventHandler() {
	config.RegisterEventHandler(common.ROOT+"."+common.LOG_LEVEL, func(es []*model.Event) {
		if len(es) > 0 {
			level := es[0].Value.(string)
			s.SetLevel(level)
		}
	})
}

func (s *logger) registerNamespaceEventHandler() {
	config.RegisterEventHandler(common.ROOT+"."+common.NODE+"."+common.APP_NAMESPACE, func(es []*model.Event) {
		if len(es) > 0 {
			ns := es[0].Value.(string)
			internal.SetNamespace(ns)
		}
	})
}

func (s *logger) registerRotateConfigEventHandler() {
	config.RegisterEventHandler(common.ROOT+"."+common.LOG_ROTATE_CONFIG, func(es []*model.Event) {
		if len(es) > 0 {
			for _, e := range es {
				value, ok := e.Value.(int)
				if !ok || value <= 0 {
					continue
				}
				switch e.Key {
				case common.ROOT + "." + common.LOG_ROTATE_CONFIG_MAXSIZE:
					s.cfg.RotateConfig.MaxSize = value
				case common.ROOT + "." + common.LOG_ROTATE_CONFIG_MAXBACKUPS:
					s.cfg.RotateConfig.MaxBackups = value
				case common.ROOT + "." + common.LOG_ROTATE_CONFIG_MAXAGE:
					s.cfg.RotateConfig.MaxAge = value
				default:
					continue
				}
			}
			s.UpdateWriteSyncs()
		}
	})
}
