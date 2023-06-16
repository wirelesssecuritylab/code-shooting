package loggersinkx

import (
	"net/url"
	"runtime"

	"code-shooting/infra/logger"
)

type LogReporter interface {
	Report(logline string)
	Flush()
	Close()
}

type ReporterFactory func(*url.URL) (LogReporter, error)

func RegisterSink(factory ReporterFactory) error {
	b := &sinkBuilder{
		reporterFactory: factory,
	}
	err := logger.RegisterSink(b)
	if err != nil {
		logger.Error("loggersinkx register sink failed:", err)
		return err
	}
	return nil
}

type sinkBuilder struct {
	reporterFactory ReporterFactory
}

func (s *sinkBuilder) Build(url *url.URL) (logger.Sink, error) {
	r, err := s.reporterFactory(url)
	if err != nil {
		return nil, err
	}
	runtime.SetFinalizer(r, func(reporter LogReporter) {
		reporter.Flush()
		reporter.Close()
	})
	return &lokiSink{
		url:      url,
		reporter: r,
	}, nil
}

func (s *sinkBuilder) Scheme() string {
	return "loki"
}

type lokiSink struct {
	url      *url.URL
	reporter LogReporter
}

func (s *lokiSink) Sync() error {
	s.reporter.Flush()
	return nil
}

func (s *lokiSink) Close() error {
	return nil
}

func (s *lokiSink) Write(p []byte) (n int, err error) {
	logline := string(p)
	s.reporter.Report(logline)
	return len(logline), nil
}
