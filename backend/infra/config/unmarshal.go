package config

import (
	"fmt"
	"reflect"
	"strings"

	"code-shooting/infra/config/internal/log"

	"code-shooting/infra/config/model"

	ghodssyaml "github.com/ghodss/yaml"
	"github.com/pkg/errors"
	"go.uber.org/config"
	conf "go.uber.org/config"
	yamlv2 "gopkg.in/yaml.v2"
)

type Unmarshaler interface {
	Unmarshal(key string, obj interface{}) error
}

func newUnmarshaler(store model.Store) Unmarshaler {
	return &unmarshaler{
		store:       store,
		configValue: make(map[string]interface{}),
	}
}

type unmarshaler struct {
	configValue map[string]interface{}
	store       model.Store
}

func (s *unmarshaler) Unmarshal(key string, obj interface{}) error {
	res, err := s.handle("")
	if err != nil {
		return err
	}
	yaml, err := conf.NewYAML(conf.Static(res), conf.Permissive())
	if err != nil {
		return err
	}
	if s.jsonInline(obj) {
		return s.unmarshalByJson(key, obj, yaml)
	}
	value := yaml.Get(key)
	if !value.HasValue() {
		return fmt.Errorf(errKeyNotExist+" in the config file: %s", key)
	}
	err = value.Populate(obj)
	if err != nil {
		return errors.Wrap(err, "get key from config file")
	}
	return nil
}

func (s *unmarshaler) handle(prefix string) (interface{}, error) {
	s.configValue = s.getAllConfigs()
	res := s.handleMap(prefix, s.configValue)
	return res, nil

}

func (s *unmarshaler) handleMap(prefix string, config map[string]interface{}) interface{} {
	result := make(map[string]interface{})
	mapKeys := s.getMapKeys(config, prefix)
	for _, key := range mapKeys {
		keyArr := strings.Split(key, "#")
		if configItem, ok := config[prefix+key]; ok {
			s.setMapIndex(result, keyArr, configItem)
		}
	}
	return result
}

func (s *unmarshaler) getAllConfigs() map[string]interface{} {
	configs := make(map[string]interface{}, 0)
	for _, item := range s.store.GetAll() {
		configs[item.Key] = item.Value
	}
	return configs
}

func (s *unmarshaler) setMapIndex(result map[string]interface{}, keySlice []string, value interface{}) {
	tmp := result
	len := len(keySlice)
	lastIndex := len - 1
	for i := 0; i < lastIndex; i++ {
		key := keySlice[i]
		var ok bool = false
		m, ok := tmp[key]
		if !ok {
			newMap := make(map[string]interface{}, 1)
			tmp[key] = newMap
			tmp = newMap
		} else {
			tmp, ok = m.(map[string]interface{})
			if !ok {
				log.Errorf("get value by key[%s] is not map", key)
				return
			}
		}
	}
	res := s.getConfigValue(value)
	tmp[keySlice[lastIndex]] = res
}

func (s *unmarshaler) getConfigValue(value interface{}) interface{} {
	switch value.(type) {
	case []interface{}:
		return s.handleSlice(value.([]interface{}))
	case map[string]interface{}:
		return s.handleMap("", value.(map[string]interface{}))
	default:
		return value
	}
}

func (s *unmarshaler) getMapKeys(configValue map[string]interface{}, prefix string) []string {
	var mapKeys []string
	pLen := len(prefix)
	for key := range configValue {
		isPrefix := checkPrefix(key, prefix)
		if !isPrefix {
			continue
		}
		mapKeys = append(mapKeys, key[pLen:])
	}
	return mapKeys
}

func (s *unmarshaler) handleSlice(configs []interface{}) interface{} {
	result := make([]interface{}, 0)
	for _, config := range configs {
		switch config.(type) {
		case []interface{}:
			result = append(result, s.handleSlice(config.([]interface{})))
		case map[string]interface{}:
			result = append(result, s.handleMap("", config.(map[string]interface{})))
		default:
			result = append(result, config)
		}

	}
	return result
}

func checkPrefix(key, prefix string) bool {
	if key != "" && prefix == "" {
		return true
	}
	// kLen := len(key)
	// pLen := len(prefix)
	// if kLen < pLen {
	// 	return false
	// }
	// if key[:pLen] == prefix && (kLen == pLen || key[pLen] == '.') {
	// 	return true
	// }
	return false
}

func (s *unmarshaler) jsonInline(configureObj interface{}) bool {
	objType := reflect.TypeOf(configureObj)
	return s.containsJsonInline(objType)
}

func (s *unmarshaler) containsJsonInline(t reflect.Type) bool {
	if s.hasElem(t) {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return false
	}

	for i := 0; i < t.NumField(); i++ {
		if s.containsJsonInline(t.Field(i).Type) {
			return true
		}

		if !t.Field(i).Anonymous {
			continue
		}
		jsonValue := t.Field(i).Tag.Get("json")
		if strings.Contains(jsonValue, "inline") {
			return true
		}

	}
	return false
}

func (s *unmarshaler) hasElem(t reflect.Type) bool {
	kind := t.Kind()
	switch kind {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Ptr, reflect.Slice:
		return true
	default:
		return false
	}
}

func (s *unmarshaler) unmarshalByJson(key string, configureObj interface{}, yaml *config.YAML) error {

	yamlContent, err := yamlv2.Marshal(yaml.Get(key).Value())
	if err != nil {
		return errors.Wrap(err, "config yaml marshal to bytes stream")
	}
	err = ghodssyaml.Unmarshal(yamlContent, configureObj)
	if err != nil {
		return errors.Wrap(err, "config bytes stream unmarshal to object")
	}

	return nil
}
