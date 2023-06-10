package workspacecfg

import (
	"encoding/json"
	"io/ioutil"

	"code-shooting/infra/logger"
)

func ReadWorkSpaceCfg(filepath string) []WorkSpaceDTO {
	var res []WorkSpaceDTO
	f, err := ioutil.ReadFile(filepath)
	if err != nil {
		logger.Errorf("read config file %s failed: ", filepath)
		return res
	}
	err = json.Unmarshal([]byte(f), &res)
	if err != nil {
		logger.Error("parse config file %s failed: ", filepath)
		return res
	}
	return res
}
