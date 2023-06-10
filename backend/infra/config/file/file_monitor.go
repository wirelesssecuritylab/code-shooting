package file

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"sync"
	"time"

	"code-shooting/infra/config/internal/log"
	"code-shooting/infra/config/internal/utils"
	"code-shooting/infra/config/model"

	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
)

const watchDebounceDelay = time.Second

func newFileMonitor(source model.ConfigSource, files []string, handle FileHandler) (*FileMonitor, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		err := w.Add(file)
		if err != nil {
			w.Close()
			return nil, errors.Wrap(err, "add file to watcher")
		}
	}
	return &FileMonitor{
		watcher: w,
		files:   files,
		source:  source,
		handle:  handle,
	}, nil
}

type FileMonitor struct {
	watcher     *fsnotify.Watcher
	files       []string
	handle      FileHandler
	source      model.ConfigSource
	eventsCache []*model.Event
	mutex       sync.RWMutex
}

func (s *FileMonitor) Start(stop <-chan struct{}) error {
	err := s.checkAndUpdateConfig()
	if err != nil {
		return err
	}
	ch := make(chan struct{})
	go s.watch(ch, stop)
	return nil
}

func (s *FileMonitor) watch(ch chan struct{}, stop <-chan struct{}, path ...string) {
	defer s.watcher.Close()
	var debounce <-chan time.Time
	for {
		select {
		case <-debounce:
			debounce = nil
			err := s.checkAndUpdateConfig()
			if err != nil {
				log.Errorf("read file err : %s", err.Error())
			}

		case ev := <-s.watcher.Events:
			if debounce == nil {
				debounce = time.After(watchDebounceDelay)
			}
			s.handleEvents(ev)
		case err := <-s.watcher.Errors:
			log.Warn("Error watching file trigger: %v %v", path, err)
			return
		case signal := <-stop:
			log.Info("Shutting down file watcher: %v %v", path, signal)
			return
		}
	}

}

func (s *FileMonitor) handleEvents(ev fsnotify.Event) {
	if ev.Op == fsnotify.Remove || ev.Op == fsnotify.Rename || ev.Op == fsnotify.Chmod {
		s.watcher.Remove(ev.Name)
		time.Sleep(500 * time.Millisecond)
		s.watcher.Add(ev.Name)
		log.Debug("ev name:%s,ev Op:%v", ev.Name, ev.Op)
	}
}

func (s *FileMonitor) readConfigFiles() ([]*model.ConfigItem, error) {
	configMap, err := s.handle(s.files)
	if err != nil {
		return nil, err
	}
	configSlice := make([]*model.ConfigItem, 0, len(configMap))
	now := time.Now()
	for k, v := range configMap {
		configSlice = append(configSlice, &model.ConfigItem{Key: k, Value: v, Version: now.String()})
	}
	return configSlice, nil
}

func (s *FileMonitor) checkAndUpdateConfig() error {
	newConfigs, err := s.readConfigFiles()
	if err != nil {
		return err
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	sort.Sort(model.ConfigItemSlice(newConfigs))
	oldConfigs := s.source.GetAll()
	sort.Sort(model.ConfigItemSlice(oldConfigs))
	oldIndex, newIndex := s.compareConfigItem(oldConfigs, newConfigs)
	oldLen := len(oldConfigs)
	newLen := len(newConfigs)
	for ; oldIndex < oldLen; oldIndex++ {
		log.Debug("Delete: %s ", oldConfigs[oldIndex].Key)
		s.delete(oldConfigs[oldIndex].Key)
	}
	for ; newIndex < newLen; newIndex++ {
		config := newConfigs[newIndex]
		log.Debug("Create: %s ", config.Key)
		s.create(config.Key, config.Value)
	}
	s.source.ProcessConfigEvent(s.eventsCache)
	s.eventsCache = make([]*model.Event, 0)
	return nil
}

func (s *FileMonitor) compareConfigItem(old, new []*model.ConfigItem) (int, int) {
	oldLen := len(old)
	newLen := len(new)
	oldIndex, newIndex := 0, 0
	for oldIndex < oldLen && newIndex < newLen {
		oldConfig := old[oldIndex]
		newConfig := new[newIndex]
		if v := strings.Compare(oldConfig.Key, newConfig.Key); v < 0 {
			log.Debug("delete: %s", oldConfig.Key)
			s.delete(oldConfig.Key)
			oldIndex++
		} else if v > 0 {
			log.Debug("Create: %s ", newConfig.Key)
			s.create(newConfig.Key, newConfig.Value)
			newIndex++
		} else {
			if !reflect.DeepEqual(oldConfig.Value, newConfig.Value) {
				log.Debug("Update: %s ", newConfig.Key)
				s.update(newConfig.Key, newConfig.Value)
			}
			oldIndex++
			newIndex++
		}
	}
	return oldIndex, newIndex
}

func (s *FileMonitor) create(key string, value interface{}) (string, error) {
	ver, err := s.source.Create(key, value)
	if err != nil {
		return "", err
	}
	item := s.source.Get(key)
	if item != nil {
		replica := *item
		replica.Key = utils.ConvertInnerKeyToOut(replica.Key)
		s.eventsCache = append(s.eventsCache, &model.Event{EventType: model.Create, ConfigItem: replica})
	}
	return ver, nil
}

func (s *FileMonitor) update(key string, value interface{}) (string, error) {
	version, err := s.source.Update(key, value)
	if err != nil {
		return "", err
	}
	item := s.source.Get(key)
	if item != nil {
		replica := *item
		replica.Key = utils.ConvertInnerKeyToOut(replica.Key)
		s.eventsCache = append(s.eventsCache, &model.Event{EventType: model.Update, ConfigItem: replica})
	}
	return version, nil
}

func (s *FileMonitor) delete(key string) error {
	item := s.source.Get(key)
	if item == nil {
		return fmt.Errorf("key: %s  not found", key)
	}
	replica := *item
	s.source.Delete(key)
	replica.Key = utils.ConvertInnerKeyToOut(replica.Key)
	s.eventsCache = append(s.eventsCache, &model.Event{EventType: model.Delete, ConfigItem: replica})
	return nil
}
