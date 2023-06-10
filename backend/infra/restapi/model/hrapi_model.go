package model

type QueryInfoInput struct {
	PageNo            int    `json:"pageNo"`
	PageSize          int    `json:"pageSize"`
	Msname            string `json:"msname"`
	QueryConditionDTO DTO    `json:"queryConditionDTO"`
}

type DTO struct {
	QueryType string `json:"queryType"`
	QueryVer  string `json:"queryVer"`
	QueryKey  string `json:"queryKey"`
}

type QueryInfoOutput struct {
	Bo   QueryBO           `json:"bo"`
	Code TokenResponseCode `json:"code"`
}

type QueryBO struct {
	Rows []QueryRow `json:"rows"`
}

type QueryRow struct {
	EmpUIID     string `json:"empUIID"`
	OrgNamePath string `json:"orgNamePath"`
	EmpName     string `json:"empName"`
	OrgId       string `json:"orgID"`
}

type GetInfoInput struct {
	Msname     string      `json:"msname"`
	IdType     string      `json:"idType"`
	Ids        []string    `json:"ids"`
	InfoBlocks []InfoBlock `json:"infoBlocks"`
}

type InfoBlock struct {
	Block string `json:"block"`
	Ver   string `json:"ver"`
}

type GetInfoOutput struct {
	Bo    map[string]StaffMessage `json:"bo"`
	Code  TokenResponseCode       `json:"code"`
	Other interface{}             `json:"other"`
}

type StaffMessage struct {
	CompanyID   string `json:"companyId"`
	CompanyName string `json:"companyName"`
	OrgID       string `json:"orgId"`
	OrgName     string `json:"orgName"`
	OrgFullName string `json:"orgFull_Name"`
	OfficeNO    string `json:"officeNo"`
	OfficeName  string `json:"officeName"`
}

type TokenResponseCode struct {
	Code  string `json:"code"`
	Msg   string `json:"msg"`
	MsgId string `json:"msgId"`
}
