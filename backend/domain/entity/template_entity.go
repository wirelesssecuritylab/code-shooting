package entity

type TemplateEntity struct {
	Id        string `json:"id"`
	Version   string `json:"version"`
	Active    bool   `json:"active"`
	Workspace string `json:"workspace"`
	UploadBy  string `json:"uploadBy"`
	UploadAt  int64  `json:"uploadAt"`
}

type TemplateOpHistory struct {
	Id             string `json:"id"`
	Action         string `json:"action"`
	CurrentVersion string `json:"currentVersion"`
	NextVersion    string `json:"nextVersion"`
	Changlog       string `json:"changlog"`
	Operator       string `json:"operator"`
	OpTime         int64  `json:"opTime"`
	OpStatus       string `json:"opStatus"`
	Workspace      string `json:"workspace"`
}
