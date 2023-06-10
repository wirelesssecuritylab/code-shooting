package entity

type DeptUserMappings struct {
	DeptUserRels []DeptUserRel `json:"mappings"`
}

type DeptUserRel struct {
	Dept         string       `json:"depart_name,omitempty"`
	UsersMapping UsersMapping `json:"mapping,omitempty"`
}

type UsersMapping struct {
	Users []string `json:"staff_id"`
}
