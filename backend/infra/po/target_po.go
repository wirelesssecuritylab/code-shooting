package po

import (
	"fmt"

	"time"

	"code-shooting/infra/database/pg/sql"

	pg "github.com/lib/pq"
	"github.com/pkg/errors"
)

type TargetPo struct {
	Id              string         `gorm:"primary_key;not null" json:"id,omitempty"`
	Name            string         `gorm:"unique;not null" json:"name,omitempty"`
	Language        string         `json:"language,omitempty"`
	Template        string         `json:"template,omitempty"`
	Owner           string         `json:"owner,omitempty"`
	OwnerName       string         `json:"ownerName,omitempty"`
	IsShared        bool           `json:"isShared,omitempty"`
	TagId           string         `json:"TagId,omitempty"`
	Answer          string         `json:"answer,omitempty"`
	Targets         pg.StringArray `gorm:"type:text[]" json:"targets"`
	CustomLabelInfo string         `json:"customLableInfo,omitempty"`
	ExtendedLabel   pg.StringArray `gorm:"type:text[]" json:"extendedLabel"`
	InstituteLabel  pg.StringArray `gorm:"type:text[]" json:"instituteLabel"`
	RelatedRanges   pg.StringArray `gorm:"type:text[]" json:"relatedRanges"`
	Workspace       string         `json:"workspace,omitempty"`

	MainCategory string    `json:"mainCategory,omitempty"`
	SubCategory  string    `json:"subCategory,omitempty"`
	DefectDetail string    `json:"defectDetail,omitempty"`
	CreateTime   time.Time `json:"createTime,omitempty"`
	UpdateTime   time.Time `json:"updateTime,omitempty"`

	TotalAnswerNum   int `json:"totalAnswerNum,omitempty"`
	TotalAnswerScore int `json:"totalAnswerScore,omitempty"`
}

type TargetDB struct {
	*sql.GormDB
}

func (m TargetDB) InsertTarget(target *TargetPo) error {
	result := m.Unscoped().Where("id=?", target.Id).Delete(target)
	if result.Error == nil {
		result = m.Model(target).Save(target)
	}
	return errors.WithStack(result.Error)
}

func (m TargetDB) UpdateTarget(target *TargetPo) error {
	result := m.Model(target).Where("id=?", target.Id).Select("*").Updates(target)
	return errors.WithStack(result.Error)
}

func (m TargetDB) DeleteTarget(target *TargetPo) error {
	result := m.Model(target).Where("id=?", target.Id).Delete(target)
	return errors.WithStack(result.Error)
}

func (m TargetDB) FindTargets(where, value string) ([]TargetPo, error) {
	targets := make([]TargetPo, 0)
	result := m.Find(&targets, fmt.Sprintf("%s=?", where), value)
	return targets, errors.WithStack(result.Error)
}

func (m TargetDB) FindTarget(value string) (TargetPo, error) {
	var t TargetPo
	result := m.Where("id = ?", value).First(&t)
	return t, result.Error
}

func (m TargetDB) IsExist(target *TargetPo) (bool, error) {
	targets := make([]TargetPo, 0)
	result := m.Model(target).Where("id=?", target.Id).Find(&targets)

	if result.RowsAffected == 0 && result.Error == nil {
		return false, nil
	}
	if result.Error != nil {
		return false, errors.WithStack(result.Error)
	}
	return true, nil
}

func (m TargetDB) FindAll() ([]TargetPo, error) {
	targets := make([]TargetPo, 0)
	result := m.Model(&TargetPo{}).Find(&targets)
	return targets, result.Error
}
