package customuser

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"sync"

	"code-shooting/infra/logger"

	"code-shooting/domain/entity"
	"code-shooting/infra/util"
)

type CustomUserService struct {
	usrDeptMap map[string]string
}

var _custUserServiceMu sync.Mutex
var custUserService *CustomUserService

func GetCustomUserService() *CustomUserService {
	_custUserServiceMu.Lock()
	defer _custUserServiceMu.Unlock()
	if custUserService == nil {
		custUserService = &CustomUserService{usrDeptMap: buildUserDeptMap(loadAllCfgDeptUsers())}
	}
	return custUserService
}
func ReloadCustomUserService(_ string) error {
	_custUserServiceMu.Lock()
	defer _custUserServiceMu.Unlock()
	custUserService = &CustomUserService{usrDeptMap: buildUserDeptMap(loadAllCfgDeptUsers())}
	return nil
}
func (s *CustomUserService) FindDeptOfCustomUser(user string) string {
	if res, ok := s.usrDeptMap[user]; !ok {
		logger.Infof("FindDeptOfCustomUser '%s' cant find '%s'", user, res)
		return res
	} else {
		return res
	}

}

func buildUserDeptMap(deptUsers *entity.DeptUserMappings) map[string]string {
	userDetpMap := map[string]string{}
	for _, dept := range deptUsers.DeptUserRels {
		for _, user := range dept.UsersMapping.Users {
			userDetpMap[user] = dept.Dept
		}
	}
	return userDetpMap
}

func loadAllCfgDeptUsers() *entity.DeptUserMappings {
	filepath := filepath.Join(util.ConfDir, "project", "deptUserMapping.json")
	f, err := ioutil.ReadFile(filepath)
	if err != nil {
		logger.Infof("read config file failed, ", err.Error())
		return &entity.DeptUserMappings{}
	}

	depts := &entity.DeptUserMappings{}
	err = json.Unmarshal([]byte(f), depts)
	if err != nil {
		logger.Infof("parse config file failed, ", err.Error())
		return &entity.DeptUserMappings{}
	}
	return depts
}
