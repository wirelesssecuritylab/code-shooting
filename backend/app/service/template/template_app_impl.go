package template

import (
	"code-shooting/domain/entity"
	"code-shooting/domain/service/template"
	"code-shooting/infra/util"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"

	"code-shooting/infra/logger"

	"github.com/google/uuid"
)

type TemplateAppServiceImpl struct {
	template.ITemplateService
}

var templateAppService ITemplateAppService

func GetTemplateAppService() ITemplateAppService {
	if templateAppService == nil {
		templateAppService = &TemplateAppServiceImpl{template.GetTemplateService()}
	}
	return templateAppService
}

func GetTemplateVerFromFileName(fileName string) (string, error) {
	re := regexp.MustCompile("^(.+)-v((\\d)+)\\.((\\d)+)\\.xlsm$")
	if ok := re.MatchString(fileName); !ok {
		return "", errors.New("file name is not right")
	}
	return fileName[strings.LastIndex(fileName, "-")+1 : len(fileName)-5], nil
}

func (impl *TemplateAppServiceImpl) UploadTemplate(uploadFile io.Reader, uploadFileName, operator, version string, workspace string) error {
	err := impl.UploadTmpFile(uploadFile, uploadFileName, workspace)
	defer impl.insertUploadOpHistory(version, operator, workspace, err)
	if err != nil {
		logger.Errorf("upload tmp file failed: %v", err.Error())
		return err
	}

	template := newTemplate(version, operator, false, workspace)

	if err = impl.InsertTemplate(template); err != nil {
		logger.Errorf("insert template failed: %v", err.Error())
		return err
	}
	if err = impl.MoveTemplateFile(uploadFileName, workspace); err != nil {
		logger.Errorf("move template file to dir failed: %v", err.Error())
		impl.RemoveTemplate(template)
		return err
	}
	return nil
}

const (
	OpActionEnable  = "enable"
	OpActionDisable = "disable"
)

func (impl *TemplateAppServiceImpl) EnableOrDisableTemplate(id string, opHistory *entity.TemplateOpHistory) error {
	targetTemplate, err := impl.QueryTemplateByVersion(opHistory.CurrentVersion, opHistory.Workspace)
	defer impl.insertOpHistory(opHistory, err)
	if err != nil {
		return err
	}

	if opHistory.Action == OpActionDisable {
		targetTemplate.Active = false
		return impl.updateTemplate(targetTemplate)
	}

	targetTemplate.Active = true
	err = impl.updateTemplate(targetTemplate)
	if err != nil {
		return err
	}

	return nil
}

func (impl *TemplateAppServiceImpl) DeleteTemplate(id string, opHistory *entity.TemplateOpHistory) error {
	temp, err := impl.QueryTemplateById(id)
	defer impl.insertOpHistory(opHistory, err)
	if err != nil {
		logger.Errorf("query template failed: %v", err.Error())
		return err
	}
	if err = impl.RemoveTemplate(temp); err != nil {
		logger.Errorf("delete template failed: %v", err.Error())
		return err
	}
	if err = impl.RemoveFile(BuildTemplateFileName(temp.Version), temp.Workspace); err != nil {
		logger.Errorf("remove file failed: %v", err.Error())
		impl.InsertTemplate(temp)
	}
	return nil
}

func (impl *TemplateAppServiceImpl) InitDefaultTemplate() error {
	templates, err := impl.QueryTemplate()
	if err != nil {
		logger.Errorf("query tempaltes failed: %v", err.Error())
		return err
	}
	if len(templates) == 0 {
		fileNames, err := impl.GetTemplateFiles(util.DefaultWorkspace)
		if err != nil {
			logger.Errorf("get tempalte files failed: %v", err.Error())
			return err
		}
		impl.initTemplateFiles(fileNames)
	}
	return nil
}

func (impl *TemplateAppServiceImpl) initTemplateFiles(fileNames []string) {
	tempStatus := true
	for _, fileName := range fileNames {
		version, err := GetTemplateVerFromFileName(fileName)
		if err != nil {
			continue
		}
		temp, err := impl.QueryTemplateByVersion(version, util.DefaultWorkspace)
		if err != nil {
			logger.Errorf("get template by version failed: %v", err.Error())
			continue
		}
		if temp == nil {
			err = impl.InsertTemplate(newTemplate(version, "admin", tempStatus, util.DefaultWorkspace))
			if err != nil {
				continue
			}
			tempStatus = false
		}
	}
}

func (impl *TemplateAppServiceImpl) updateTemplate(template *entity.TemplateEntity) error {
	if err := impl.UpdateTempate(template); err != nil {
		logger.Errorf("update templates failed: %v", err.Error())
		return err
	}
	return nil
}

func (impl *TemplateAppServiceImpl) insertUploadOpHistory(version, operator, workspace string, err error) {
	if err != nil {
		impl.InsertTemplateOpHistory(newUploadOpHistory(version, operator, "failed", workspace))
	} else {
		impl.InsertTemplateOpHistory(newUploadOpHistory(version, operator, "success", workspace))
	}
}

func (impl *TemplateAppServiceImpl) insertOpHistory(opHistory *entity.TemplateOpHistory, err error) {
	if err != nil {
		impl.InsertTemplateOpHistory(fillOpHistory(opHistory, "failed"))
	} else {
		impl.InsertTemplateOpHistory(fillOpHistory(opHistory, "success"))
	}
}

func newUploadOpHistory(version, operator, opStatus string, workspace string) *entity.TemplateOpHistory {
	return &entity.TemplateOpHistory{
		Id:             uuid.NewString(),
		Action:         "add",
		CurrentVersion: version,
		NextVersion:    "",
		Changlog:       "上传规范",
		Operator:       operator,
		OpTime:         time.Now().Unix(),
		OpStatus:       opStatus,
		Workspace:      workspace,
	}
}

func fillOpHistory(entity *entity.TemplateOpHistory, status string) *entity.TemplateOpHistory {
	entity.Id = uuid.NewString()
	entity.OpTime = time.Now().Unix()
	entity.OpStatus = status
	return entity
}

func newTemplate(version, operator string, active bool, workspace string) *entity.TemplateEntity {
	return &entity.TemplateEntity{
		Id:        uuid.NewString(),
		Version:   version,
		Active:    active,
		Workspace: workspace,
		UploadBy:  operator,
		UploadAt:  time.Now().Unix(),
	}
}

func BuildTemplateFileName(version string) string {
	return fmt.Sprintf("代码打靶落地模板-%s.xlsm", version)
}
