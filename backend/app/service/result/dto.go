package result

type ResultRespDto struct {
	UserId     string            `json:"userId,omitempty"`
	UserName   string            `json:"userName,omitempty"`
	Department string            `json:"department,omitempty"`
	TeamName   string            `json:"teamName,omitempty"`
	CenterName string            `json:"centerName,omitempty"`
	Targets    []TargetResultDto `json:"targets,omitempty"`
	RangeScore RangeLangScoreDto `json:"rangeScore,omitempty"`
}

type TargetResultDto struct {
	TargetId          string            `json:"targetId"`
	HitNum            int               `json:"hitNum"`
	HitScore          int               `json:"hitScore"`
	HitScoreHundredth int               `json:"hitScoreHundredth"`
	HitRate           string            `json:"hitRate"`
	TotalNum          int               `json:"totalNum"`
	TotalScore        int               `json:"totalScore"`
	Details           []TargetDetailDto `json:"detail,omitempty"`
}

type TargetDetailDto struct {
	FileName       string `json:"fileName"`
	StartLineNum   int    `json:"startLineNum"`
	EndLineNum     int    `json:"endLineNum"`
	StartColNum    int    `json:"startColNum"`
	EndColNum      int    `json:"endColNum"`
	DefectClass    string `json:"defectClass"`
	DefectSubClass string `json:"defectSubClass"`
	DefectDescribe string `json:"defectDescribe"`
	TargetScore    int    `json:"score"`
	Remark         string `json:"remark"`
}

type RangeLangScoreDto struct {
	HitNum            int    `json:"hitNum"`
	HitScore          int    `json:"hitScore"`
	HitScoreHundredth int    `json:"hitScoreHundredth"`
	HitRate           string `json:"hitRate"`
	TotalNum          int    `json:"totalNum"`
	TotalScore        int    `json:"totalScore"`
}

type TargetID struct {
	TargetID string `json:"targetId"`
}
type Targets struct {
	Targets []TargetID `json:"targets"`
}
