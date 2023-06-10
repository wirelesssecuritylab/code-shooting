package staff_agg

import (
	"code-shooting/domain/dto"
	role_agg "code-shooting/domain/privilege/role-agg"
	"code-shooting/infra/util/tools"
	"fmt"
)

func NewStaffEntitys(staffs dto.StaffRolesDto, roles dto.RolePrivilegesDto) ([]StaffEntity, error) {
	res := make([]StaffEntity, 0, len(staffs.Mappings))
	allRoleIds := role_agg.ParseRoleMapToRoleIds(roles)
	var err error
	for _, staff := range staffs.Mappings {
		var staffRoleIds []string
		for _, roleId := range staff.Mapping.RoleIds {
			if !tools.IsContain(allRoleIds, roleId) {
				err = fmt.Errorf("allroleids '%s' not contian  '%v' \n %v", allRoleIds, roleId, err)
				continue
			}
			staffRoleIds = append(staffRoleIds, roleId)
		}
		res = append(res, StaffEntity{Id: staff.Id, RoleIdVos: staffRoleIds})
	}
	return res, err
}
