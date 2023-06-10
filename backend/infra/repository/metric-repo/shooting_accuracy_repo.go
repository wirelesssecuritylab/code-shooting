package metricrepository

import (
	metricData "code-shooting/infra/po/metric-data-po"
	"code-shooting/infra/util/database"
)

type ShootingAccuracyRepository struct {
	ShootingAccuracyDB metricData.ShootingAccuracyDB
}

func NewShootingAccuracyRepository() *ShootingAccuracyRepository {
	return &ShootingAccuracyRepository{
		ShootingAccuracyDB: metricData.ShootingAccuracyDB{GormDB: database.DB},
	}
}

func (s *ShootingAccuracyRepository) Save(e *metricData.ShootingAccuracyPo) error {
	return s.ShootingAccuracyDB.SaveShootingAccuracy(e)
}

func (s *ShootingAccuracyRepository) SaveInBatch(e *[]metricData.ShootingAccuracyPo) error {
	return s.ShootingAccuracyDB.SaveShootingAccuracyInBatch(e)
}
