package po

import (
	"time"

	pg "github.com/lib/pq"
)

type TbEC struct {
	Id                string         `gorm:"primary_key;not null"`
	Institute         string         `gorm:"not null"`
	Center            string         `gorm:"not null"`
	Department        string         `gorm:"not null"`
	Team              string         `gorm:"not null"`
	Importer          string         `gorm:"not null"`
	ImportTime        time.Time      `gorm:"not null"`
	ConvertedToTarget bool           `gorm:"not null"`
	AssociatedTargets pg.StringArray `gorm:"type:text[];not null"`
}

func (s *TbEC) TableName() string {
	return "ec"
}
