package model

type WXCMQueryInput struct {
	AppKey          string           `json:"appKey"`
	AppSecret       string           `json:"appSecret"`
	AppVersion      string           `json:"appVersion"`
	ModelName       string           `json:"modelName"`
	EqualConditions []EqualCondition `json:"equalCondition"`
}

type EqualCondition struct {
	Column string `json:"column"`
	Value  string `json:"value"`
}

type WXCMQueryOutput struct {
	Bo   WXCMBO                `json:"bo"`
	Code WXCMTokenResponseCode `json:"code"`
}

type WXCMBO struct {
	Columns     []string  `json:"columns"`
	HasNextPage int       `json:"hasNextPage"`
	Rows        []WXCMRow `json:"rows"`
}

type WXCMRow struct {
	DepartmentCode     string `json:"department_code"`
	DepartmentName     string `json:"department_name"`
	InstituteCode      string `json:"institute_code"`
	RdCenterCode       string `json:"rd_center_code"`
	RdCenterName       string `json:"rd_center_name"`
	EmployeeName       string `json:"employee_name"`
	InstituteName      string `json:"institute_name"`
	EmployeeCode       string `json:"employee_code"`
	EmployeeCodeName   string `json:"employee_code_name"`
	ProductSystemName  string `json:"product_system_name"`
	ProductSystemCode  string `json:"product_system_code"`
	Level2DivisionName string `json:"level2_division_name"`
	Level2DivisionCode string `json:"level2_division_code"`
	PostName           string `json:"post_name"`
	CompanyName        string `json:"company_name"`
	CompanyCode        string `json:"company_code"`
	EmployeeStatus     string `json:"employee_status"`
}

type WXCMTokenResponseCode struct {
	Code  string `json:"code"`
	Msg   string `json:"msg"`
	MsgId string `json:"msgId"`
}
