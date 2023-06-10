package entity

import "time"

type Range struct {
	Id         string
	Name       string
	Type       string
	ProjectId  string
	Owner      string
	Targets    []string
	StartTime  time.Time
	EndTime    time.Time
	CreateTime time.Time
	UpdateTime time.Time
	DesensTime time.Time
}
