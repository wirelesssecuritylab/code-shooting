package privilegecfg

import (
	"code-shooting/domain/dto"
	"encoding/json"
	"io/ioutil"

	"code-shooting/infra/logger"
)

func StaffRoleCfgParseRead(filepath string) (dto.StaffRolesDto, error) {
	var res = dto.StaffRolesDto{}
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

type StaffRoleCfgParse func(string) (dto.StaffRolesDto, error)

func (f StaffRoleCfgParse) Read(filepath string) (dto.StaffRolesDto, error) {
	return f(filepath)
}
