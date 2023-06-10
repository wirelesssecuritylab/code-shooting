package dto

type RangeShootingResult struct {
	UserName string               `json:"userName"`
	UserId   string               `json:"userId"`
	Targets  []SubmitTargetResult `json:"targets"`
}

type SubmitTargetResult struct {
	TargetId       string `json:"targetId"`
	FileName       string `json:"fileName"`
	StartLineNum   int    `json:"startLineNum"`
	EndLineNum     int    `json:"endLineNum"`
	StartColNum    int    `json:"startColNum"`
	EndColNum      int    `json:"endColNum"`
	DefectClass    string `json:"defectClass"`
	DefectSubClass string `json:"defectSubClass"`
	DefectDescribe string `json:"defectDescribe"`
	Remark         string `json:"remark"`
}

type TargetAnswers struct {
	TargetID string          `json:"targetid"`
	Answers  []*TargetAnswer `json:"answers"`
}

type TargetAnswer struct {
	FileName       string `json:"fileName"`
	StartLineNum   int    `json:"startLineNum"`
	EndLineNum     int    `json:"endLineNum"`
	DefectClass    string `json:"defectClass"`
	DefectSubClass string `json:"defectSubClass"`
	DefectDescribe string `json:"defectDescribe"`
}
