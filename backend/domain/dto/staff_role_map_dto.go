package dto

type Roles struct {
	RoleIds []string `json:"role_ids"`
}

type StaffRolesDto struct {
	Mappings []StaffRolesMapping `json:"mappings"`
}

type StaffRolesMapping struct {
	Id      string `json:"id"`
	Mapping Roles  `json:"mapping"`
}
