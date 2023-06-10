package model

type Store interface {
	Get(key string) *ConfigItem
	Create(key string, value interface{}) (string, error)
	Update(key string, value interface{}) (string, error)
	Delete(key string) error
	GetAll() []*ConfigItem
}
