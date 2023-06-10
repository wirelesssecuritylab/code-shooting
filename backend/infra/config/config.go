package config

import (
	"strings"

	"code-shooting/infra/config/model"
)

type EventHandler model.EventHandler

type Config interface {
	Get(key string, value interface{}) error
	GetValue(key string) interface{}
}

const errKeyNotExist = "the key does not exist"

func IsNotExist(err error) bool {
	return strings.Contains(err.Error(), errKeyNotExist)
}
