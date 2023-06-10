package assembler

import (
	"code-shooting/infra/logger"
	"encoding/json"
)

func ParseReq(body []byte, req interface{}) error {
	err := json.Unmarshal(body, req)
	if err != nil {
		logger.Errorf("ParseReq json.Unmarshal %s err %s", string(body), err.Error())
		return err
	}
	//logger.Debugf("ParseReq result:%#v", req)
	return nil
}
