package po

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type StringSlice []string

func (s *StringSlice) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return errors.New("value not string")
	}
	return json.Unmarshal([]byte(str), s)
}

func (s StringSlice) Value() (driver.Value, error) {
	bytes, err := json.Marshal(s)
	return string(bytes), err
}

func (s StringSlice) GormDataType() string {
	return "text"
}

type TbRange struct {
	Id         string      `gorm:"primary_key;not null"`
	Name       string      `gorm:"unique;not null"`
	Type       string      `gorm:"not null"`
	ProjectId  string      `gorm:"not null"`
	Owner      string      `gorm:"not null"`
	Targets    StringSlice `gorm:"not null"`
	StartTime  time.Time   `gorm:"not null"`
	EndTime    time.Time   `gorm:"not null"`
	CreateTime time.Time   `gorm:"not null"`
	UpdateTime time.Time   `gorm:"not null"`
	DesensTime time.Time
}

func (s *TbRange) TableName() string {
	return "range"
}
