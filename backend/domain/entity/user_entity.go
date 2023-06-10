package entity

type UserEntity struct {
	Id         string   `json:"id"`
	Name       string   `json:"name"`
	Department string   `json:"department"`
	OrgId      string   `json:"orgID"`
	Institute  string   `json:"institute"`
	Email      string   `json:"email"`
	TeamName   string   `json:"teamName"`
	CenterName string 	`json:"centerName"`
	Privileges []string `json:"privileges"`
}

type QrCodeEntity struct {
	QrCodeKey   string `json:"qrCodeKey"`
	QrCodeValue string `json:"qrCodeValue"`
}
