package dto

type ShootingNoteDto struct {
	UserName string              `json:"username"`
	UserId   string              `json:"userid"`
	TargetId string              `json:"targetid"`
	RangeID  string              `json:"rangeid"`
	Targets  []ShootingResultDto `json:"targets"`
}

type ShootingResultDto struct {
	SubmitTargetResult
	ScoreNum int `json:"scorenum"`
}

type ShootingDraftDto struct {
	UserName string               `json:"username"`
	UserId   string               `json:"userid"`
	TargetId string               `json:"targetid"`
	RangeID  string               `json:"rangeid"`
	Targets  []SubmitTargetResult `json:"targets"`
}
