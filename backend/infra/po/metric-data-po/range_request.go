package metricpo

import (
	"time"

	"code-shooting/infra/database/pg/sql"
)

type RangeRequestPo struct {
	ID     string `gorm:"primary_key;not null" json:"id"`
	UserID string `json:"userid"`

	RequestTime time.Time `json:"requesttime"`
}

type RangeRequestDB struct {
	*sql.GormDB
}

func (m *RangeRequestDB) SaveRangeRequest(a *RangeRequestPo) error {
	result := m.Model(a).Save(a)
	return result.Error
}

func (m *RangeRequestDB) TableName() string {
	return "range_request"
}
