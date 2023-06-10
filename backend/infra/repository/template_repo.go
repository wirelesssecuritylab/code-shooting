package repository

import (
	"code-shooting/domain/entity"
	"code-shooting/infra/po"
	"code-shooting/infra/util/database"

	"code-shooting/infra/database/pg/sql"
	"code-shooting/infra/logger"

	"github.com/pkg/errors"
)

type TemplateRepositoryImpl struct {
	db *sql.GormDB
}

func NewTemplateRepository() *TemplateRepositoryImpl {
	return &TemplateRepositoryImpl{db: database.DB}
}

func (impl *TemplateRepositoryImpl) InsertTemplate(entity *entity.TemplateEntity) error {
	if entity == nil {
		return errors.New("tempate is null")
	}
	tempPo := templateEntity2Po(entity)
	result := impl.db.Unscoped().Where("id=?", tempPo.Id).Delete(tempPo)
	if result.Error == nil {
		result = impl.db.Model(tempPo).Save(tempPo)
	}
	return errors.WithStack(result.Error)
}

func (impl *TemplateRepositoryImpl) InsertTemplateOpHistory(entity *entity.TemplateOpHistory) {
	if entity == nil {
		return
	}
	templateHis := templateOpHistoryEntity2Po(entity)
	result := impl.db.Unscoped().Where("id = ?", templateHis.Id).Delete(templateHis)
	if result.Error == nil {
		result = impl.db.Model(templateHis).Save(templateHis)
	}
	if result.Error != nil {
		logger.Errorf("insert template operation history failed: %v", result.Error.Error())
	}
}

func (impl *TemplateRepositoryImpl) RemoveTemplate(entity *entity.TemplateEntity) error {
	tempPo := templateEntity2Po(entity)
	result := impl.db.Where("id = ?", tempPo.Id).Delete(tempPo)
	return errors.WithStack(result.Error)
}

func (impl *TemplateRepositoryImpl) UpdateTempate(template *entity.TemplateEntity) error {
	templatePo := templateEntity2Po(template)
	result := impl.db.Model(templatePo).Where("id = ?", templatePo.Id).Select("*").Updates(templatePo)
	if result.Error != nil {
		return errors.WithStack(result.Error)
	}
	return nil
}
func (impl *TemplateRepositoryImpl) UpdateActiveFalseByWorksapce(workspace string) error {
	result := impl.db.Model(&po.TemplatePo{}).Where("workspace = ?", workspace).Update("active", false)
	if result.Error != nil {
		return errors.WithStack(result.Error)
	}
	return nil
}
func (impl *TemplateRepositoryImpl) IsWorkSpaceEnable(workspace string) (bool, error) {
	tempPos := make([]po.TemplatePo, 0)
	err := impl.db.Where("active = ?", true).Where("workspace = ?", workspace).Find(&tempPos).Error
	if err != nil {
		return false, errors.WithStack(err)
	}
	if len(tempPos) < 1 {
		return false, nil
	}
	return true, nil
}
func (impl *TemplateRepositoryImpl) QueryTemplateById(id string) (*entity.TemplateEntity, error) {
	tempPo := &po.TemplatePo{}
	result := impl.db.Where("id = ?", id).First(tempPo)
	if result.Error != nil {
		return nil, errors.WithStack(result.Error)
	}
	return templatePo2Entity(tempPo), nil
}

func (impl *TemplateRepositoryImpl) QueryTemplateByVersion(version string, workspace string) (*entity.TemplateEntity, error) {
	tempPos := make([]po.TemplatePo, 0)
	result := impl.db.Where("version = ?", version).Where("workspace = ?", workspace).Find(&tempPos)
	if result.Error != nil {
		return nil, errors.WithStack(result.Error)
	}
	if len(tempPos) > 0 {
		return templatePo2Entity(&tempPos[0]), nil
	}
	return nil, nil
}

func (impl *TemplateRepositoryImpl) QueryTemplate() ([]entity.TemplateEntity, error) {
	tempPos := make([]po.TemplatePo, 0)
	result := impl.db.Find(&tempPos).Order("uploadAt DESC")
	if result.Error != nil {
		return nil, errors.WithStack(result.Error)
	}
	temps := make([]entity.TemplateEntity, 0, len(tempPos))
	for index := range tempPos {
		temps = append(temps, *templatePo2Entity(&tempPos[index]))
	}
	return temps, nil
}

func (impl *TemplateRepositoryImpl) QueryTemplateOpHistory() ([]entity.TemplateOpHistory, error) {
	tempOpHisPos := make([]po.TemplateOpHistoryPo, 0)
	result := impl.db.Find(&tempOpHisPos).Order("opTime DESC")
	if result.Error != nil {
		return nil, errors.WithStack(result.Error)
	}
	tempOpHis := make([]entity.TemplateOpHistory, 0, len(tempOpHisPos))
	for index := range tempOpHisPos {
		tempOpHis = append(tempOpHis, *templateOpHistoryPo2Entity(&tempOpHisPos[index]))
	}
	return tempOpHis, nil
}

func templateEntity2Po(entity *entity.TemplateEntity) *po.TemplatePo {
	return &po.TemplatePo{
		Id:        entity.Id,
		Version:   entity.Version,
		Active:    entity.Active,
		Workspace: entity.Workspace,
		UploadBy:  entity.UploadBy,
		UploadAt:  entity.UploadAt,
	}
}

func templatePo2Entity(tempPo *po.TemplatePo) *entity.TemplateEntity {
	return &entity.TemplateEntity{
		Id:        tempPo.Id,
		Version:   tempPo.Version,
		Active:    tempPo.Active,
		Workspace: tempPo.Workspace,
		UploadBy:  tempPo.UploadBy,
		UploadAt:  tempPo.UploadAt,
	}
}

func templateOpHistoryEntity2Po(entity *entity.TemplateOpHistory) *po.TemplateOpHistoryPo {
	return &po.TemplateOpHistoryPo{
		Id:             entity.Id,
		Action:         entity.Action,
		CurrentVersion: entity.CurrentVersion,
		NextVersion:    entity.NextVersion,
		Changlog:       entity.Changlog,
		Operator:       entity.Operator,
		OpTime:         entity.OpTime,
		OpStatus:       entity.OpStatus,
		Workspace:      entity.Workspace,
	}
}

func templateOpHistoryPo2Entity(tempOpHisPo *po.TemplateOpHistoryPo) *entity.TemplateOpHistory {
	return &entity.TemplateOpHistory{
		Id:             tempOpHisPo.Id,
		Action:         tempOpHisPo.Action,
		CurrentVersion: tempOpHisPo.CurrentVersion,
		NextVersion:    tempOpHisPo.NextVersion,
		Changlog:       tempOpHisPo.Changlog,
		Operator:       tempOpHisPo.Operator,
		OpTime:         tempOpHisPo.OpTime,
		OpStatus:       tempOpHisPo.OpStatus,
		Workspace:      tempOpHisPo.Workspace,
	}
}
