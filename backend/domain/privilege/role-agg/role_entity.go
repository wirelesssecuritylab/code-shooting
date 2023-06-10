package role_agg

type RoleEntity struct {
	Id           string
	Name         string
	PrivilegeVos []string
}

func (s RoleEntity) GetPrivileges() []string {
	return s.PrivilegeVos
}
