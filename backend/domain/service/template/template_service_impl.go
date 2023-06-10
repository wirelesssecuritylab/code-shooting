package template

import (
	"io"
	"path/filepath"

	repo "code-shooting/domain/repository"
	"code-shooting/infra/repository"
	"code-shooting/infra/util"
	"code-shooting/infra/util/tools"
)

func templateTmpDir() string {
	return filepath.Join(util.ConfDir, ".template_tmp")
}

func templateDir() string {
	return filepath.Join(util.ConfDir, "templates")
}

type TemplateServiceImpl struct {
	repo.ITemplateRepository
}

var templateService ITemplateService

func GetTemplateService() ITemplateService {
	if templateService == nil {
		templateService = &TemplateServiceImpl{repository.NewTemplateRepository()}
	}
	return templateService
}

func (impl *TemplateServiceImpl) UploadTmpFile(uploadFile io.Reader, uploadFileName string, workspace string) error {
	return tools.SaveFile(filepath.Join(templateTmpDir(), workspace, uploadFileName), uploadFile)
}

func (impl *TemplateServiceImpl) RemoveFile(fileName string, workspace string) error {
	return tools.RemoveFile(filepath.Join(templateDir(), workspace, fileName))
}

func (impl *TemplateServiceImpl) MoveTemplateFile(fileName string, workspace string) error {
	return tools.MoveFile(filepath.Join(templateTmpDir(), workspace, fileName), filepath.Join(templateDir(), workspace, fileName))
}

func (impl *TemplateServiceImpl) GetTemplateFiles(workspace string) ([]string, error) {
	return tools.ListFileNames(filepath.Join(templateDir(), workspace))
}
