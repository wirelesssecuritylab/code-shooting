package role_agg

import (
	"code-shooting/domain/dto"
	"code-shooting/infra/util/tools"
	"fmt"
)

func NewRoleEntitys(roles dto.RolePrivilegesDto, privileges dto.PrivilegesDto) ([]RoleEntity, error) {
	res := make([]RoleEntity, 0, len(roles.Mappings))
	allRoleIds := ParseRoleMapToRoleIds(roles)
	var err error
	allPrivileges, err := parsePrivilegesDto(privileges)
	if err != nil {
		return nil, err
	}
	for _, role := range roles.Mappings {
		var rolePs []string
		for _, roleP := range role.Mapping.PrivilegeNames {
			if !tools.IsContain(allPrivileges, roleP) {
				err = fmt.Errorf("privilege %s not contian in %v \n %v", roleP, allPrivileges, err)
				continue
			}
			rolePs = append(rolePs, roleP)
		}
		if !tools.IsContain(allRoleIds, role.Id) {
			err = fmt.Errorf("role id %s not in all roleids %v", role.Id, allRoleIds)
			continue
		}
		res = append(res, RoleEntity{Id: role.Id, Name: role.Name, PrivilegeVos: rolePs})
	}
	return res, err
}

func parsePrivilegesDto(privileges dto.PrivilegesDto) ([]string, error) {
	var res []string
	var err error
	for _, layer := range privileges.Mappings {
		for _, privilege := range layer.Mapping.PrivilegeNames {
			if tools.IsContain(res, privilege) {
				err = fmt.Errorf("privilege %s has contian in %v \n %v", privilege, res, err)
				continue
			}
			res = append(res, privilege)
		}
	}
	return res, err
}
func ParseRoleMapToRoleIds(roles dto.RolePrivilegesDto) []string {
	var res []string
	for _, role := range roles.Mappings {
		res = append(res, role.Id)
	}
	return res
}
