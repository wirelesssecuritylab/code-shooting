package config

import (
	"code-shooting/infra/config/file"
	"code-shooting/infra/config/internal/log"
	"code-shooting/infra/config/model"
	"reflect"
	"runtime"

	"github.com/pkg/errors"
)

var aggSource aggConfigSource

func init() {
	aggSource = newAggregateConfigSource()
}

func NewConfig(configPath string, opts ...Option) (Config, error) {
	param := &configOptionParam{}
	Options(opts).Do(param)
	fileSource, err := file.NewFileSource(file.DefaultYamlSource, file.ParseYamlFile, configPath)
	if err != nil {
		return nil, err
	}
	if !param.fileChangeMonitor {
		return newConfig(fileSource)
	}
	err = AddConfigSource(fileSource)
	if err != nil {
		return nil, err
	}
	return newConfig(aggSource)
}

func newConfig(source model.ConfigSource) (Config, error) {
	config := &refactorConfig{source: source, stop: make(chan struct{})}
	err := config.Start()
	if err != nil {
		return nil, err
	}
	return config, nil
}

type refactorConfig struct {
	source model.ConfigSource
	stop   chan struct{}
}

func (s *refactorConfig) Get(key string, obj interface{}) error {
	rv := reflect.ValueOf(obj)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.New("invalid object supplied")
	}
	um := newUnmarshaler(s.source)
	return um.Unmarshal(key, obj)
}

func (s *refactorConfig) GetValue(key string) interface{} {
	item := s.source.Get(key)
	if item == nil {
		return nil
	}
	return item.Value
}

func (s *refactorConfig) Start() error {
	runtime.SetFinalizer(s, func(config *refactorConfig) {
		s.stop <- struct{}{}
	})
	return s.source.Start(s.stop)
}

func AddConfigSource(source model.ConfigSource) error {
	return aggSource.AddSource(source)
}

func RegisterEventHandler(key string, handler EventHandler) {
	aggSource.RegisterEventHandler(key, model.EventHandler(handler))
}

func ProcessConfigEvent(events []*model.Event) {
	aggSource.ProcessConfigEvent(events)
}

func SetLogger(l log.Logger) {
	log.SetLogger(l)
}
