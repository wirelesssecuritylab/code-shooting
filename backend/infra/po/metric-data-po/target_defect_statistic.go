package metricpo

import (
	"code-shooting/infra/database/pg/sql"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type TargetDefectStatPo struct {
	ID        string `gorm:"primary_key;not null" json:"id"`
	TargetID  string
	DefectId  string
	DefectNum int
}

type TargetDefectStatDB struct {
	*sql.GormDB
}

func (m *TargetDefectStatDB) TableName() string {
	return "target_defect_statistic"
}

func (m *TargetDefectStatDB) SaveTargetDefectStat(a *TargetDefectStatPo) error {
	result := m.Model(a).Save(a)
	return result.Error
}

func (m *TargetDefectStatDB) SaveTargetDefectStatInBatch(a *[]TargetDefectStatPo) error {
	result := m.Model(a).CreateInBatches(a, len(*a))
	return result.Error
}

func (m *TargetDefectStatDB) UpdateTargetDefectStatInBatch(a *[]TargetDefectStatPo) error {
	if len(*a) == 0 {
		return nil
	}
	result := m.Transaction(func(tx *gorm.DB) error {
		if err := m.Delete((*a)[0].TargetID); err != nil {
			return errors.WithStack(err)
		}
		return m.SaveTargetDefectStatInBatch(a)
	})
	return result
}

func (m *TargetDefectStatDB) Delete(targetid string) error {
	result := m.Model(&TargetDefectStatPo{}).Where("target_id=?", targetid).Delete(&TargetDefectStatPo{})
	return result.Error
}

func (m *TargetDefectStatDB) DeleteAll() error {
	result := m.Model(&TargetDefectStatPo{}).Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&TargetDefectStatPo{})
	return result.Error
}
