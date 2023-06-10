package po

type TbResult struct {
	Id         int64  `gorm:"primary_key;auto_increment"`
	UserId     string `gorm:"not null"`
	RangeId    string `gorm:"not null"`
	TargetId   string `gorm:"not null"`
	HitNum     int    `gorm:"not null"`
	TotalNum   int    `gorm:"not null"`
	HitScore   int    `gorm:"not null"`
	TotalScore int    `gorm:"not null"`
	Language   string `gorm:"not null"` // lower case in db
	Records    []byte `gorm:"not null"`
}

func (s *TbResult) TableName() string {
	return "result"
}
