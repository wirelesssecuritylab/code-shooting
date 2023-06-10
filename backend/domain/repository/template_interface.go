package repository

import "code-shooting/domain/entity"

type ITemplateRepository interface {
	InsertTemplate(entity *entity.TemplateEntity) error
	InsertTemplateOpHistory(entity *entity.TemplateOpHistory)
	RemoveTemplate(entity *entity.TemplateEntity) error
	UpdateTempate(template *entity.TemplateEntity) error
	UpdateActiveFalseByWorksapce(workspace string) error
	IsWorkSpaceEnable(workspace string) (bool, error)
	QueryTemplateById(id string) (*entity.TemplateEntity, error)
	QueryTemplateByVersion(version string, workspace string) (*entity.TemplateEntity, error)
	QueryTemplate() ([]entity.TemplateEntity, error)
	QueryTemplateOpHistory() ([]entity.TemplateOpHistory, error)
}
