package workspace

import (
	"path"

	"code-shooting/infra/util"
	"code-shooting/infra/workspacecfg"
)

var cfgPath = "workspace/workspace.json"

type WorkSpaceAppService struct {
}

func GetWorkSpaceAppService() *WorkSpaceAppService {
	return &WorkSpaceAppService{}
}

func (w *WorkSpaceAppService) GetWorkSpace() []workspacecfg.WorkSpaceDTO {
	return workspacecfg.ReadWorkSpaceCfg(path.Join(util.ConfDir, cfgPath))
}
