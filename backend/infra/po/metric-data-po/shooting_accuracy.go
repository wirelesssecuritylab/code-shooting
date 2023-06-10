package metricpo

import (
	"time"

	"code-shooting/infra/database/pg/sql"
)

type ShootingAccuracyPo struct {
	ID         string `gorm:"primary_key;not null" json:"id"`
	UserID     string `json:"userid"`
	RangeID    string `json:"rangeid"`
	TargetID   string `json:"targetid"`
	DefectType string `json:"defectType"`
	HitNum     int    `json:"hitNum"`
	HitScore   int    `json:"hitScore"`

	SubmitTime time.Time `json:"submitTime"`
}

type ShootingAccuracyDB struct {
	*sql.GormDB
}

func (m *ShootingAccuracyDB) SaveShootingAccuracy(a *ShootingAccuracyPo) error {
	result := m.Model(a).Save(a)
	return result.Error
}

func (m *ShootingAccuracyDB) SaveShootingAccuracyInBatch(a *[]ShootingAccuracyPo) error {
	result := m.Model(a).CreateInBatches(a, len(*a))
	return result.Error
}

func (m *ShootingAccuracyDB) TableName() string {
	return "shooting_accuracy"
}
