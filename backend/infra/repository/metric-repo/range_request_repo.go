package metricrepository

import (
	metricData "code-shooting/infra/po/metric-data-po"
	"code-shooting/infra/util/database"
)

type RangeRequestRepository struct {
	RangeRequestDB metricData.RangeRequestDB
}

func NewRangeRequestRepository() *RangeRequestRepository {
	return &RangeRequestRepository{
		RangeRequestDB: metricData.RangeRequestDB{GormDB: database.DB},
	}
}

func (s *RangeRequestRepository) Save(e *metricData.RangeRequestPo) error {
	return s.RangeRequestDB.SaveRangeRequest(e)
}
