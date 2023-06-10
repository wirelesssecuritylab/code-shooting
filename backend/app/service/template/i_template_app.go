package template

import (
	"code-shooting/domain/entity"
	"io"
)

type ITemplateAppService interface {
	QueryTemplateByVersion(version string, workspace string) (*entity.TemplateEntity, error)
	UploadTemplate(uploadFile io.Reader, uploadFileName, operator, version string, workspace string) error
	EnableOrDisableTemplate(id string, entity *entity.TemplateOpHistory) error
	DeleteTemplate(id string, entity *entity.TemplateOpHistory) error
	QueryTemplate() ([]entity.TemplateEntity, error)
	QueryTemplateOpHistory() ([]entity.TemplateOpHistory, error)
	InitDefaultTemplate() error
}
