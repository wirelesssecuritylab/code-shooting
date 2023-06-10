package dto

type TemplateAction struct {
	Action string        `json:"name"`
	Params TemplateModel `json:"parameters"`
}

type TemplateModel struct {
	TempleteId     string `json:"templateId"`
	Action         string `json:"action"`
	CurrentVersion string `json:"currentVersion"`
	NextVersion    string `json:"nextVersion"`
	Changlog       string `json:"changlog"`
	Operator       string `json:"operator"`
	Worksapce      string `json:"workspace"`
}

type DefectDetail struct {
	Description string `json:"description"`
	Code        string `json:"code"`
}

type Defects map[string]map[string][]DefectDetail
