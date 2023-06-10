package dto

type IdModel struct {
	Id string `json:"id"`
}

type QrCodeModel struct {
	QrCodeKey   string `json:"qrCodeKey"`
	QrCodeValue string `json:"qrCodeValue"`
}

type TargetAction struct {
	Action string      `json:"name"`
	Target TargetModel `json:"parameters"`
}

type TargetModel struct {
	ID             string      `json:"id"`
	Name           string      `json:"name"`
	User           string      `json:"user"`
	Owner          string      `json:"owner"`
	OwnerName      string      `json:"ownerName"`
	Lang           string      `json:"language"`
	Template       string      `json:"template"`
	TagID          string      `json:"tagId"`
	RangeID        string      `json:"rangeid"`
	TagName        TagNameInfo `json:"tagName"`
	CustomLabel    string      `json:"customLable"`
	ExtendedLabel  []string    `json:"extendedLabel"`
	InstituteLabel []string    `json:"instituteLabel"`
	Answer         string      `json:"answer"`
	Isshared       bool        `json:"isShared"`
	Files          []string    `json:"targets"`
	RelatedRanges  []string    `json:"relatedRanges"`
	Workspace      string      `json:"workspace"`
}

func (t *TargetModel) IsQueryUserValid() bool {
	if len(t.Owner) == 0 && len(t.User) == 0 {
		return false
	}
	return true
}

func (t *TargetModel) QueryUserIsTargetOwner(targetOwner string) bool {
	if t.Owner == targetOwner || t.User == targetOwner {
		return true
	}

	return false
}

type TagNameInfo struct {
	MainCategory string `json:"mainCategory"`
	SubCategory  string `json:"subCategory"`
	DefectDetail string `json:"defectDetail"`
}

type DefectAction struct {
	Lang            string `json:"language"`
	NeedCode        bool   `json:"needCode"`
	TemplateVersion string `json:"templateVersion"`
}
type TargetCodeFile struct {
	Filename string `json:"filename"`
}
