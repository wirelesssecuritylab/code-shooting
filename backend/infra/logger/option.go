package logger

import (
	"code-shooting/infra/logger/internal"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Option func(interface{})

type Options []Option

func (s Options) Do(i interface{}) {
	for _, option := range s {
		option(i)
	}
}

// Please don't rely on the declare
// Recommend use encapsulated TimeEncoder/CustomFormat interface

func TimeEncoder(layout string) Option {
	return func(i interface{}) {
		cfg, ok := i.(*zapcore.EncoderConfig)
		if ok {
			cfg.EncodeTime = internal.TimeEncoderOfLayout(layout)
		}
	}
}

func CapitalLevelEncoder(levelMapping map[string]string) Option {
	return func(i interface{}) {
		cfg, ok := i.(*zapcore.EncoderConfig)
		if ok {
			mapping := make(map[string]string)
			for k, v := range levelMapping {
				mapping[strings.ToUpper(k)] = strings.ToUpper(v)
			}
			cfg.EncodeLevel = func(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
				level, ok := mapping[l.CapitalString()]
				if ok {
					enc.AppendString(level)
				} else {
					enc.AppendString(l.CapitalString())
				}
			}
		}
	}
}

func CustomFormat(format string) Option {
	return func(i interface{}) {
		cfg, ok := i.(*configDto)
		if ok {
			cfg.Format = format
		}
	}
}

func AddCallerSkip(skip int) Option {
	return func(i interface{}) {
		zapOptions, ok := i.(*[]zap.Option)
		if ok {
			*zapOptions = append(*zapOptions, zap.AddCallerSkip(skip))
		}
	}
}
