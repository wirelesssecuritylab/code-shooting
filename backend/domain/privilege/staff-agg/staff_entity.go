package staff_agg

type StaffEntity struct {
	Id        string
	RoleIdVos []string
}

func (s StaffEntity) GetRoleIds() []string {
	return s.RoleIdVos
}
