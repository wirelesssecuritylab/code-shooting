package target

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"code-shooting/app/service/template"
	"code-shooting/domain/service/target"
	shootingresult "code-shooting/infra/shooting-result"
	"code-shooting/infra/shooting-result/defect"
	"code-shooting/infra/util"
)

type TargetAppService struct {
}

func GetTargetAppService() *TargetAppService {
	return &TargetAppService{}
}

func (t *TargetAppService) CheckTarget(id string) ([]TargetAnswerDTO, error) {
	var res []TargetAnswerDTO

	tg, err := target.GetTargetService().FindTarget(id)
	if err != nil {
		return nil, err
	}
	answerFile := filepath.Join(tg.GetAnswerFileDir(), tg.Answer)

	answers, err := shootingresult.NewShootingResultCalculator().LoadShootingData(answerFile)
	if err != nil {
		return nil, err
	}

	fileName, _ := getCurrentTemplateFileName(tg.Workspace, tg.Template)
	coder, err := defect.NewDefectEncoder(filepath.Join(util.TemplateDir, tg.Workspace, fileName))
	if err != nil {
		return nil, err
	}
	for _, a := range answers {
		defectClass, _, _ := coder.DecodeDefect(a.DefectCode)
		fmt.Println(defectClass)
		if defectClass == "" {
			res = append(res, TargetAnswerDTO{a.FileName, a.StartLineNum, a.EndLineNum, a.DefectClass, a.DefectSubClass, a.DefectDescribe})
		}
	}
	return res, nil
}

type TargetAnswerDTO struct {
	FileName     string `json:"fileName"`
	StartLineNum int    `json:"startLineNum"`
	EndLineNum   int    `json:"endLineNum"`
	Class        string `json:"class"`
	SubClass     string `json:"subClass"`
	Describe     string `json:"describe"`
}

func getCurrentTemplateFileName(workspace, version string) (string, error) {
	if workspace == "" {
		workspace = util.DefaultWorkspace
	}
	templates, err := template.GetTemplateAppService().QueryTemplate()
	if err != nil {
		return "", err
	}
	if version != "" {
		for index := range templates {
			if templates[index].Active && templates[index].Workspace == workspace && templates[index].Version == version {
				return template.BuildTemplateFileName(templates[index].Version), nil
			}
		}
	} else {
		var versionsActive []string
		for index := range templates {
			if templates[index].Active && templates[index].Workspace == workspace {
				versionsActive = append(versionsActive, (templates[index].Version)[1:len((templates[index].Version))])
			}
		}
		if len(versionsActive) == 0 {
			return "", fmt.Errorf("工作空间 %s 中不存在启用中规范", workspace)
		}
		maxVersion := versionsActive[0]
		maxVersion_a, _ := strconv.Atoi(strings.Split(maxVersion, ".")[0])
		maxVersion_b, _ := strconv.Atoi(strings.Split(maxVersion, ".")[1])
		for i := 0; i < len(versionsActive); i++ {
			version_a, _ := strconv.Atoi(strings.Split(versionsActive[i], ".")[0])
			version_b, _ := strconv.Atoi(strings.Split(versionsActive[i], ".")[1])
			if maxVersion_a < version_a {
				maxVersion_a = version_a
				maxVersion_b = version_b
			} else if maxVersion_a == version_a {
				if maxVersion_b < version_b {
					maxVersion_b = version_b
				}
			}
		}
		return template.BuildTemplateFileName("v" + strconv.Itoa(maxVersion_a) + "." + strconv.Itoa(maxVersion_b)), nil
	}
	return "", fmt.Errorf("工作空间 %s 中不存在启用中规范", workspace)
}
