package monitor

import (
	"strings"
	"sync"

	"code-shooting/infra/config/model"
)

func NewConfigstoreMonitor() model.Monitor {
	return &configStoreMonitor{
		handlers: make(map[string][]model.EventHandler),
	}
}

type configStoreMonitor struct {
	handlers map[string][]model.EventHandler
	mutex    sync.RWMutex
}

func (s *configStoreMonitor) Start(stop <-chan struct{}) error {
	return nil
}

func (s *configStoreMonitor) ProcessConfigEvent(events []*model.Event) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	for key, handles := range s.handlers {
		var handleEvents []*model.Event
		for _, e := range events {
			if strings.HasPrefix(e.Key, key) || key == "" {
				handleEvents = append(handleEvents, e)
			}
		}
		if len(handleEvents) == 0 {
			continue
		}
		for _, handle := range handles {
			handle(handleEvents)
		}
	}
}

func (s *configStoreMonitor) RegisterEventHandler(key string, handler model.EventHandler) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	handlers, _ := s.handlers[key]
	handlers = append(handlers, handler)
	s.handlers[key] = handlers
}
