package project

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"sync"

	"code-shooting/infra/logger"

	"code-shooting/domain/entity"
	customuser "code-shooting/domain/service/custom-user"
	"code-shooting/domain/service/user"
	"code-shooting/infra/util"
)

const DefaultProjectId = "public"
const DefaultProjectName = "公共"

type ProjectService struct {
	projects   map[string]*entity.Project
	deptPrjMap map[string][]string
}

var _prjServiceMu sync.Mutex
var prjService *ProjectService

func GetProjectService() *ProjectService {
	_prjServiceMu.Lock()
	defer _prjServiceMu.Unlock()
	if prjService == nil {
		prjService = loadProjectService()
	}
	return prjService
}

func ReloadProjectService(_ string) error {
	_prjServiceMu.Lock()
	defer _prjServiceMu.Unlock()
	prjService = loadProjectService()
	return nil
}
func loadProjectService() *ProjectService {
	prjService = &ProjectService{projects: make(map[string]*entity.Project)}
	projsCfg := loadAllProjects()
	for _, p := range projsCfg.Projects {
		prjService.projects[p.Id] = p
	}
	prjService.projects[DefaultProjectId] = &entity.Project{Id: DefaultProjectId, Name: DefaultProjectName}
	prjService.deptPrjMap = buildDeptPrjMap(projsCfg)
	return prjService
}

func (s *ProjectService) FindByUser(user *entity.UserEntity) []string {
	return s.findPrjsByDept(user.Department, customuser.GetCustomUserService().FindDeptOfCustomUser(user.Id))
}

func (s *ProjectService) FindByUserId(id string) []*entity.Project {
	var res []*entity.Project
	user, err := user.NewUserDomainService().QueryUser(&entity.UserEntity{Id: id})
	if err != nil {
		return nil
	}
	projectIds := s.findPrjsByDept(user.Department)
	for _, projectId := range projectIds {
		p := s.FindPrjsByID(projectId)
		if p == nil {
			continue
		}
		res = append(res, p)
	}
	return res
}

func (s *ProjectService) findPrjsByDept(depts ...string) []string {
	results := []string{DefaultProjectId}
	for _, dept := range depts {
		prjs, exist := s.deptPrjMap[dept]
		if !exist {
			continue
		}
		results = append(results, prjs...)
	}
	return results
}

func (s *ProjectService) FindPrjsByID(id string) *entity.Project {
	if prj, exist := s.projects[id]; exist {
		result := *prj
		return &result
	}
	return nil
}

func (s *ProjectService) ListProjects() []*entity.Project {
	projs := make([]*entity.Project, 0, len(s.projects))
	for _, p := range s.projects {
		projs = append(projs, p)
	}
	return projs
}

func buildDeptPrjMap(prjs *entity.ProjectsMappings) map[string][]string {
	results := map[string][]string{}
	for _, prj := range prjs.Projects {
		for _, dept := range prj.DeptsMapping.Depts {
			if _, exist := results[dept]; !exist {
				results[dept] = []string{prj.Id}
			} else {
				results[dept] = append(results[dept], prj.Id)
			}
		}
	}
	logger.Infof("buildDeptPrjMap = %v", results)
	return results
}

func loadAllProjects() *entity.ProjectsMappings {
	filepath := filepath.Join(util.ConfDir, "privilege", "centerMaps.json")
	f, err := ioutil.ReadFile(filepath)
	if err != nil {
		logger.Infof("read config file failed, ", err.Error())
		return &entity.ProjectsMappings{}
	}

	projects := &entity.ProjectsMappings{}
	err = json.Unmarshal([]byte(f), projects)
	if err != nil {
		logger.Infof("parse config file failed, ", err.Error())
		return &entity.ProjectsMappings{}
	}
	return projects
}

func (d *ProjectService) FindDeptsByProjectId(id string) []string {
	p, exist := d.projects[id]
	if !exist {
		return []string{}
	}
	return p.DeptsMapping.Depts
}
