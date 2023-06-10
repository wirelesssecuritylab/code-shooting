package config

import (
	"go.uber.org/fx"
)

var configPath = ""

func NewModule(confPath string) fx.Option {
	return fx.Provide(func(lc fx.Lifecycle) (Config, error) {
		configPath = confPath
		return NewConfig(confPath, WithFileChangeMonitor())
	})
}

func GetConfPath() string {
	return configPath
}
