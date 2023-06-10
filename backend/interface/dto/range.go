package dto

type RangeAction struct {
	Action string `json:"name"`
	Params Range  `json:"parameters"`
}

type TargetInRange struct {
	Language   string `json:"language"`
	TargetName string `json:"targetName"`
	TargetId   string `json:"targetId"`
}

type Range struct {
	Id         string          `json:"id"`
	Name       string          `json:"name"`
	Project    string          `json:"project"`
	Owner      string          `json:"owner"`
	User       string          `json:"user"`
	Type       string          `json:"type"`
	StartTime  int64           `json:"startTime"`
	EndTime    int64           `json:"endTime"`
	Targets    []TargetInRange `json:"targets"`
	DesensTime int64           `json:"desensitiveTime"`
}

func (s *Range) GetTargetIds() []string {
	var res = make([]string, 0, len(s.Targets))
	for _, t := range s.Targets {
		res = append(res, t.TargetId)
	}
	return res
}

type RangeDetail struct {
	Range
	OrgID     string `json:"orgID"`
	PrjName   string `json:"projectName"`
	OwnerName string `json:"ownerName"`
}
