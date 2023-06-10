package metricrepository

import (
	metricData "code-shooting/infra/po/metric-data-po"
	"code-shooting/infra/util/database"
)

type RingNumRepository struct {
	RingNumDB metricData.RingNumDB
}

func NewRingNumRepository() *RingNumRepository {
	return &RingNumRepository{
		RingNumDB: metricData.RingNumDB{GormDB: database.DB},
	}
}

func (s *RingNumRepository) Save(e *metricData.RingNumPo) error {
	return s.RingNumDB.SaveRingNum(e)
}

func (s *RingNumRepository) SaveInBatch(e *[]metricData.RingNumPo) error {
	return s.RingNumDB.SaveRingNumInBatch(e)
}
func (s *RingNumRepository) UpdateInBatch(e *[]metricData.RingNumPo) error {
	return s.RingNumDB.UpdateRingNumInBatch(e)
}

func (s *RingNumRepository) Update(e *metricData.RingNumPo) error {
	return s.RingNumDB.UpdateRingNum(e)
}
