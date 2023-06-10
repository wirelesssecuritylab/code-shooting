package file

import (
	"errors"

	"code-shooting/infra/config/internal/memory"
	"code-shooting/infra/config/internal/monitor"
	"code-shooting/infra/config/model"
)

const DefaultYamlSource = "YamlSource"

func NewFileSource(sourceName string, handle FileHandler, filePaths ...string) (model.ConfigSource, error) {
	if len(filePaths) <= 0 {
		return nil, errors.New("config filePaths is empty")
	}
	s := &source{
		Store:   memory.NewMemoryStore(),
		Monitor: monitor.NewConfigstoreMonitor(),
		name:    sourceName,
	}
	fileMonitor, err := newFileMonitor(s, filePaths, handle)
	if err != nil {
		return nil, err
	}
	s.fileMonitor = fileMonitor
	return s, nil
}

type source struct {
	model.Store
	model.Monitor
	fileMonitor *FileMonitor
	name        string
}

func (s *source) GetSourceName() string {
	return s.name
}

func (s *source) Start(stop <-chan struct{}) error {
	return s.fileMonitor.Start(stop)
}
