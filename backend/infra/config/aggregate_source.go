package config

import (
	"sync"

	"code-shooting/infra/config/internal/monitor"
	"code-shooting/infra/config/model"
)

type aggConfigSource interface {
	model.ConfigSource
	AddSource(model.ConfigSource) error
}

func newAggregateConfigSource() aggConfigSource {
	return &aggregateConfigSource{monitor: monitor.NewConfigstoreMonitor(), stop: make(chan struct{})}
}

type aggregateConfigSource struct {
	monitor model.Monitor
	sources sync.Map
	stop    chan struct{}
}

func (s *aggregateConfigSource) AddSource(source model.ConfigSource) error {
	source.RegisterEventHandler("", s.ProcessConfigEvent)
	s.sources.Store(source.GetSourceName(), source)
	return source.Start(s.stop)
}

func (s *aggregateConfigSource) Start(stop <-chan struct{}) error {
	s.monitor.Start(stop)
	return nil
}

func (s *aggregateConfigSource) GetSourceName() string {
	return "aggregation"
}

func (s *aggregateConfigSource) RegisterEventHandler(key string, handler model.EventHandler) {
	s.monitor.RegisterEventHandler(key, handler)
}

func (s *aggregateConfigSource) ProcessConfigEvent(events []*model.Event) {
	s.monitor.ProcessConfigEvent(events)

}

func (s *aggregateConfigSource) Get(key string) *model.ConfigItem {
	var configItem *model.ConfigItem
	s.sources.Range(func(k, v interface{}) bool {
		source := v.(model.ConfigSource)
		if item := source.Get(key); item != nil && (configItem == nil || item.Version > configItem.Version) {
			configItem = item
		}
		return true
	})
	return configItem
}

func (s *aggregateConfigSource) Create(key string, value interface{}) (string, error) {
	return "", nil
}

func (s *aggregateConfigSource) Update(key string, value interface{}) (string, error) {
	return "", nil
}

func (s *aggregateConfigSource) Delete(key string) error {
	return nil
}

func (s *aggregateConfigSource) GetAll() []*model.ConfigItem {
	configMap := make(map[string]*model.ConfigItem)
	s.sources.Range(func(k, v interface{}) bool {
		source := v.(model.ConfigSource)
		for _, configitem := range source.GetAll() {
			if configitem == nil {
				continue
			}
			v, ok := configMap[configitem.Key]
			if (!ok) || configitem.Version > v.Version {
				configMap[configitem.Key] = configitem
			}
		}
		return true
	})
	var out []*model.ConfigItem
	for _, v := range configMap {
		out = append(out, v)
	}
	return out
}
