package po

import (
	"code-shooting/infra/database/pg/sql"

	"gorm.io/gorm"
)

type DefectPo struct {
	DefectId       string `gorm:"primary_key;not null" json:"defect_code"`
	Language       string
	DefectClass    string
	DefectSubclass string
	DefectDescribe string
}

type DefectDB struct {
	*sql.GormDB
}

func (m *DefectDB) TableName() string {
	return "defect"
}

func (m *DefectDB) SaveDefect(a *DefectPo) error {
	result := m.Model(a).Save(a)
	return result.Error
}

func (m *DefectDB) SaveDefectInBatch(a *[]DefectPo) error {
	result := m.Model(a).CreateInBatches(a, len(*a))
	return result.Error
}

func (m *DefectDB) DeleteAll() error {
	result := m.Model(&DefectPo{}).Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&DefectPo{})
	return result.Error
}
