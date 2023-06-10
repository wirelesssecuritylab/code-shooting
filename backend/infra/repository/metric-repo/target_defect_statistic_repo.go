package metricrepository

import (
	metricData "code-shooting/infra/po/metric-data-po"
	"code-shooting/infra/util/database"
)

type TargetDefectStatRepository struct {
	TargetDefectStatDB metricData.TargetDefectStatDB
}

func NewTargetDefectStatRepository() *TargetDefectStatRepository {
	return &TargetDefectStatRepository{
		TargetDefectStatDB: metricData.TargetDefectStatDB{GormDB: database.DB},
	}
}

func (s *TargetDefectStatRepository) SaveInBatch(e *[]metricData.TargetDefectStatPo) error {
	return s.TargetDefectStatDB.SaveTargetDefectStatInBatch(e)
}

func (s *TargetDefectStatRepository) UpdateInBatch(e *[]metricData.TargetDefectStatPo) error {
	return s.TargetDefectStatDB.UpdateTargetDefectStatInBatch(e)
}

func (s *TargetDefectStatRepository) Remove(targetid string) error {
	return s.TargetDefectStatDB.Delete(targetid)
}

func (s *TargetDefectStatRepository) RemoveAll() error {
	return s.TargetDefectStatDB.DeleteAll()
}
