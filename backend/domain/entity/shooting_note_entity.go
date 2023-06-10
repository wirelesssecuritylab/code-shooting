package entity

import (
	"time"
)

type ShootingNoteEntity struct {
	UserID     string
	TargetID   string
	UserName   string
	RangeID    string
	Datas      []ShootingData
	UpdateTime time.Time
}

type ShootingData struct {
	FileName     string
	StartLineNum int
	EndLineNum   int
	StartColNum  int
	EndColNum    int
	Remark       string
	DefectCode   string
	ScoreNum     int
}
