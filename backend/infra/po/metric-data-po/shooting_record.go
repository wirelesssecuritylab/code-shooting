package metricpo

import (
	"time"

	"code-shooting/infra/database/pg/sql"
)

type ShootingRecordPo struct {
	ID       string `gorm:"primary_key;not null" json:"id"`
	UserID   string `json:"userid"`
	RangeID  string `json:"rangeid"`
	Language string `json:"language"`

	ShootingTime time.Time `json:"shootingTime"`
}

type ShootingRecordDB struct {
	*sql.GormDB
}

func (m *ShootingRecordDB) SaveShootingRecord(a *ShootingRecordPo) error {
	result := m.Model(a).Save(a)
	return result.Error
}

func (m *ShootingRecordDB) TableName() string {
	return "shooting_record"
}
