package metricpo

import (
	"time"

	"code-shooting/infra/database/pg/sql"
)

type ShootingDurationPo struct {
	ID       string `gorm:"primary_key;not null" json:"id"`
	UserID   string `json:"userid"`
	RangeID  string `json:"rangeid"`
	TargetID string `json:"targetid"`

	EndTime time.Time `json:"endtime"`
	Timelen int       `json:"timelen"`
}

type ShootingDurationDB struct {
	*sql.GormDB
}

func (m *ShootingDurationDB) SaveShootingDuration(a *ShootingDurationPo) error {

	result := m.Model(a).Save(a)

	return result.Error
}

func (m ShootingDurationDB) UpdateShootingDuration(a *ShootingDurationPo) error {

	shootDurationEntry := make([]*ShootingDurationPo, 0)

	m.Model(a).Where("target_id=? and user_id=? and range_id=?", a.TargetID, a.UserID, a.RangeID).Find(&shootDurationEntry)

	if int(0) != len(shootDurationEntry) {
		a.ID = shootDurationEntry[0].ID
		a.Timelen = shootDurationEntry[0].Timelen + a.Timelen
	}
	result := m.SaveShootingDuration(a)
	return result
}
func (m *ShootingDurationDB) TableName() string {
	return "shooting_Duration"
}
