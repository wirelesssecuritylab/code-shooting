package privilegecfg

import (
	"code-shooting/domain/dto"
	"encoding/json"
	"io/ioutil"

	"code-shooting/infra/logger"
)

func PrivilegeCfgParseRead(filepath string) (dto.PrivilegesDto, error) {
	var res = dto.PrivilegesDto{}
	f, err := ioutil.ReadFile(filepath)
	if err != nil {
		logger.Errorf("read config file %s failed", filepath)
		return res, err
	}
	err = json.Unmarshal([]byte(f), &res)
	if err != nil {
		logger.Error("parse config file %s failed", filepath)
		return res, err
	}
	return res, nil
}

type PrivilegeCfgParse func(string) (dto.PrivilegesDto, error)

func (f PrivilegeCfgParse) Read(filepath string) (dto.PrivilegesDto, error) {
	return f(filepath)
}
