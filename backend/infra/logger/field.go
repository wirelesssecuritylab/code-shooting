package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

// Please don't rely on the declare
// Recommend use encapsulated StringField/IntField/DurationField interface
type Field func() zapcore.Field

func StringField(k, v string) Field {
	return func() zapcore.Field {
		return zap.String(k, v)
	}
}

func IntField(k string, v int64) Field {
	return func() zapcore.Field {
		return zap.Int64(k, v)
	}
}

func DurationField(k string, v time.Duration) Field {
	return func() zapcore.Field {
		return zap.Duration(k, v)
	}
}
