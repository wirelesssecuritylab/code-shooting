package entity

type ShootingResult struct {
	UserName string
	UserId   string
	Targets  []TargetResult
}

type TargetResult struct {
	HitNum        int
	HitScore      int
	TotalNum      int
	TotalScore    int
	TargetId      string
	TargetDetails []TargetDetail
}

type TargetDetail struct {
	FileName     string
	StartLineNum int
	EndLineNum   int
	StartColNum  int
	EndColNum    int
	DefectCode   string
	TargetScore  int
	Remark       string
}
