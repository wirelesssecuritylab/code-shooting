package entity

type ProjectsMappings struct {
	Projects []*Project `json:"mappings"`
}

type Project struct {
	Id           string            `json:"id,omitempty"`
	Name         string            `json:"name,omitempty"`
	DeptsMapping DepartmentMapping `json:"mapping"`
}

type DepartmentMapping struct {
	Depts []string `json:"departments"`
}
