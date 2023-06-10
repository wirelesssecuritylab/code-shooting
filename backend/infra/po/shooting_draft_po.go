package po

import (
	"time"

	"code-shooting/infra/database/pg/sql"
	"code-shooting/infra/logger"
)

type ShootingDraftPo struct {
	ID       string `gorm:"primary_key;not null" json:"id"`
	UserID   string `json:"userid"`
	UserName string `json:"username"`
	TargetID string `json:"targetid"`
	RangeID  string `json:"rangeid"`
	Records  []byte `json:"records"`

	UpdateTime time.Time `json:"updatedtime"`
}

type ShootingDraftDB struct {
	*sql.GormDB
}

func (m *ShootingDraftDB) Save(a *ShootingDraftPo) error {
	if rec, err := m.Get(a.ID); err == nil && rec != nil {
		result := m.Model(a).Where("id=?", a.ID).Updates(a)
		logger.Infof("ShootingDraftDB Save updates, %s records:%d, result:[v%]", a.ID, len(a.Records), result.Error)
		return result.Error
	}
	result := m.Model(a).Save(a)
	logger.Infof("ShootingDraftDB Save new, %s records:%d, result:[v%]", a.ID, len(a.Records), result.Error)
	return result.Error
}

func (m *ShootingDraftDB) Update(a *ShootingDraftPo) error {
	result := m.Model(a).Where("id=?", a.ID).Updates(a)
	return result.Error
}

func (m *ShootingDraftDB) Delete(id string) error {
	result := m.Model(&ShootingDraftPo{}).Where("id=?", id).Delete(&ShootingDraftPo{})
	return result.Error
}

func (m *ShootingDraftDB) Get(id string) (*ShootingDraftPo, error) {
	u := ShootingDraftPo{}
	result := m.Find(&u, "id=?", id)

	if result.Error != nil {
		return nil, result.Error
	}

	if len(u.ID) != 0 {
		return &u, nil
	}
	return nil, nil
}

func (m *ShootingDraftDB) TableName() string {
	return "ShootingDraft"
}
