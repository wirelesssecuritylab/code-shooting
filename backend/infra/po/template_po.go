package po

type TemplatePo struct {
	Id        string `gorm:"primary_key;not null" json:"id,omitempty"`
	Version   string `gorm:"not null" json:"version,omitempty"`
	Active    bool   `json:"active,omitempty"`
	Workspace string `json:"workspace,omitempty"`
	UploadBy  string `json:"uploadBy,omitempty"`
	UploadAt  int64  `json:"uploadAt,omitempty"`
}

type TemplateOpHistoryPo struct {
	Id             string `gorm:"primary_key;not null" json:"id,omitempty"`
	Action         string `json:"action,omitempty"`
	CurrentVersion string `json:"currentVersion,omitempty"`
	NextVersion    string `json:"nextVersion,omitempty"`
	Changlog       string `json:"changelog,omitempty"`
	Operator       string `json:"operator,omitempty"`
	OpTime         int64  `json:"opTime,omitempty"`
	OpStatus       string `json:"opStatus,omitempty"`
	Workspace      string `json:"workspace,omitempty"`
}

func (s *TemplatePo) TableName() string {
	return "template"
}

func (s *TemplateOpHistoryPo) TableName() string {
	return "template_op_history"
}
