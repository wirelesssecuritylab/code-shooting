package metricrepository

import (
	metricData "code-shooting/infra/po/metric-data-po"
	"code-shooting/infra/util/database"
)

type ShootingDurationRepository struct {
	ShootingDurationDB metricData.ShootingDurationDB
}

func NewShootingDurationRepository() *ShootingDurationRepository {
	return &ShootingDurationRepository{
		ShootingDurationDB: metricData.ShootingDurationDB{GormDB: database.DB},
	}
}

func (s *ShootingDurationRepository) Save(e *metricData.ShootingDurationPo) error {
	return s.ShootingDurationDB.SaveShootingDuration(e)
}

func (s *ShootingDurationRepository) Update(e *metricData.ShootingDurationPo) error {

	return s.ShootingDurationDB.UpdateShootingDuration(e)
}
