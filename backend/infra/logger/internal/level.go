package internal

import (
	"fmt"
	"github.com/pkg/errors"
	"go.uber.org/zap/zapcore"
	"strings"
)

type Level string

const (
	DebugLevel Level = "debug"
	InfoLevel  Level = "info"
	WarnLevel  Level = "warn"
	ErrorLevel Level = "error"
	PanicLevel Level = "panic"
	FatalLevel Level = "fatal"
)

var _levelMappingTable = map[Level]zapcore.Level{
	DebugLevel: zapcore.DebugLevel,
	InfoLevel:  zapcore.InfoLevel,
	WarnLevel:  zapcore.WarnLevel,
	ErrorLevel: zapcore.ErrorLevel,
	PanicLevel: zapcore.PanicLevel,
	FatalLevel: zapcore.FatalLevel,
}

func (s *Level) Check() error {
	for level := range _levelMappingTable {
		if strings.ToLower(string(*s)) == strings.ToLower(string(level)) {
			return nil
		}
	}
	return errors.Errorf("unsupported level: %s", string(*s))
}

func (s *Level) Equal(l Level) bool {
	if err := s.Check(); err != nil {
		return false
	}
	return strings.ToLower(string(*s)) == strings.ToLower(string(l))
}

func (s *Level) ToZapLevel() (zapcore.Level, error) {
	for k, v := range _levelMappingTable {
		if s.Equal(k) {
			return v, nil
		}
	}
	return zapcore.Level(0), errors.Errorf("level %v incompatible with zap level", s.String())
}

func (s *Level) String() string {
	if err := s.Check(); err != nil {
		return fmt.Sprintf("Level(%s)", *s)
	}
	return strings.ToLower(string(*s))
}
