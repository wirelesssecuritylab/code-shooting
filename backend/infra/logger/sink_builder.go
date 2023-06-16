package logger

import (
	"net/url"
	"strings"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Sink interface {
	Sync() error
	Close() error
	Write(p []byte) (n int, err error)
}

type SinkBuilder interface {
	Build(url *url.URL) (Sink, error)
	Scheme() string
}

func RegisterSink(s SinkBuilder) error {
	err := zap.RegisterSink(s.Scheme(), func(url *url.URL) (zap.Sink, error) {
		return s.Build(url)
	})
	if err != nil && !strings.Contains(err.Error(), "already registered for scheme") {
		return err
	}
	l := GetLogger()
	if l == nil {
		return errors.New("logger is nil")
	}
	return l.UpdateWriteSyncs()
}
