package metricpo

import (
	"time"

	"code-shooting/infra/database/pg/sql"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type RingNumPo struct {
	ID         string `gorm:"primary_key;not null" json:"id"`
	OwnerID    string `json:"ownerid"`
	TargetID   string `json:"targetid"`
	DefectType string `json:"defectType"`
	RingNum    int    `json:"ringNum"`
	RingScore  int    `json:"ringScore"`

	CreateTime time.Time `json:"createTime"`
}

type RingNumDB struct {
	*sql.GormDB
}

func (m *RingNumDB) SaveRingNum(a *RingNumPo) error {
	result := m.Model(a).Save(a)
	return result.Error
}

func (m *RingNumDB) SaveRingNumInBatch(a *[]RingNumPo) error {
	result := m.Model(a).CreateInBatches(a, len(*a))
	return result.Error
}

func (m *RingNumDB) UpdateRingNumInBatch(a *[]RingNumPo) error {
	if len(*a) == 0 {
		return nil
	}
	result := m.Transaction(func(tx *gorm.DB) error {
		if err := m.Delete((*a)[0].TargetID); err != nil {
			return errors.WithStack(err)
		}
		return m.SaveRingNumInBatch(a)
	})
	return result
}

func (m *RingNumDB) Delete(targetId string) error {
	result := m.Model(&RingNumPo{}).Where("target_id=?", targetId).Delete(&RingNumDB{})
	return result.Error
}

func (m RingNumDB) UpdateRingNum(a *RingNumPo) error {
	result := m.Model(a).Where("target_id=? and defect_type=?", a.TargetID, a.DefectType).Updates(map[string]interface{}{
		"ring_num":   a.RingNum,
		"ring_score": a.RingScore,
	})
	return result.Error
}

func (m *RingNumDB) TableName() string {
	return "ring_num"
}
