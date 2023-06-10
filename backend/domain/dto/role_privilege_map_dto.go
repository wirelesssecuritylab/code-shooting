package dto

type Privileges struct {
	PrivilegeNames []string `json:"privilege_names"`
}

type RolePrivilegesDto struct {
	Mappings []RolePrivilegesMapping `json:"mappings"`
}

type PrivilegesDto RolePrivilegesDto

type RolePrivilegesMapping struct {
	Id            string     `json:"id"`
	Name          string     `json:"name"`
	BelongedLayer string     `json:"belonged_layer"`
	Mapping       Privileges `json:"mapping"`
}
