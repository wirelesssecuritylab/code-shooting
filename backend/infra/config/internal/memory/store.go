package memory

import (
	"fmt"
	"sync"
	"time"

	"code-shooting/infra/config/internal/utils"
	"code-shooting/infra/config/model"

	"github.com/labstack/gommon/log"
)

func NewMemoryStore() model.Store {
	return &memoryStore{cache: make(map[string]*model.ConfigItem)}
}

type memoryStore struct {
	cache map[string]*model.ConfigItem
	mutex sync.RWMutex
}

func (s *memoryStore) Get(key string) *model.ConfigItem {
	k := utils.ConvertOutKeyToInner(key)
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	v, ok := s.cache[k]
	if !ok {
		log.Errorf("not found by key: %s", k)
		return nil
	}
	return v
}

func (s *memoryStore) Create(key string, value interface{}) (string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	now := time.Now().String()
	_, ok := s.cache[key]
	if ok {
		return "", fmt.Errorf("key: %s ,already exists", key)
	}
	s.cache[key] = &model.ConfigItem{Key: key, Value: value, Version: now}
	return now, nil
}

func (s *memoryStore) Update(key string, value interface{}) (string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	_, ok := s.cache[key]
	if !ok {
		return "", fmt.Errorf("key: %s  not found", key)
	}
	ver := time.Now().String()
	s.cache[key] = &model.ConfigItem{Key: key, Value: value, Version: ver}
	return ver, nil
}

func (s *memoryStore) Delete(key string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.cache, key)
	return nil
}

func (s *memoryStore) GetAll() []*model.ConfigItem {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	var out []*model.ConfigItem
	for _, v := range s.cache {
		out = append(out, v)
	}
	return out
}
