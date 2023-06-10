package middleware

import (
	"code-shooting/infra/logger"

	yml "gopkg.in/yaml.v2"
)

const COMMON_ORDER int = 100

func transformInterfaceToObject(input interface{}, output interface{}) error {

	v, err := yml.Marshal(input)
	if err != nil {
		logger.Infof("Masshal error:  %v", err.Error())
		return err
	}
	err = yml.Unmarshal(v, output)
	if err != nil {
		logger.Infof("Unmasshal error:  %v", err.Error())
		return err
	}
	return nil

}
