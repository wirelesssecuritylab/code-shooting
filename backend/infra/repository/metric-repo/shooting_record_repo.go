package metricrepository

import (
	metricData "code-shooting/infra/po/metric-data-po"
	"code-shooting/infra/util/database"
)

type ShootingRecordRepository struct {
	ShootingRecordDB metricData.ShootingRecordDB
}

func NewShootingRecordRepository() *ShootingRecordRepository {
	return &ShootingRecordRepository{
		ShootingRecordDB: metricData.ShootingRecordDB{GormDB: database.DB},
	}
}

func (s *ShootingRecordRepository) Save(e *metricData.ShootingRecordPo) error {
	return s.ShootingRecordDB.SaveShootingRecord(e)
}
