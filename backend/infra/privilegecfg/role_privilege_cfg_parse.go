package privilegecfg

import (
	"code-shooting/domain/dto"
	"encoding/json"
	"io/ioutil"

	"code-shooting/infra/logger"
)

func RolePrivilegeCfgParseRead(filepath string) (dto.RolePrivilegesDto, error) {
	var res = dto.RolePrivilegesDto{}
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

type RolePrivilegeCfgParse func(string) (dto.RolePrivilegesDto, error)

func (f RolePrivilegeCfgParse) Read(filepath string) (dto.RolePrivilegesDto, error) {
	return f(filepath)
}
