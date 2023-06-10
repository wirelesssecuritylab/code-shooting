package privilegeapp

import (
	"path/filepath"

	"code-shooting/infra/logger"

	"code-shooting/domain/privilege"
	"code-shooting/infra/privilegecfg"
	"code-shooting/infra/repository"
)

func LoadPrivilegeCfg(privilegeBasePath string) error {
	var privilegesFilePath = filepath.Join(privilegeBasePath, "privileges.json")
	var rolePrivilegeMapFilePath = filepath.Join(privilegeBasePath, "role-privilege-map.json")
	var userRoleMapFilePath = filepath.Join(privilegeBasePath, "user-role-map.json")
	res := privilege.LoadPrivileges(
		privilegesFilePath, rolePrivilegeMapFilePath, userRoleMapFilePath,
		privilegecfg.PrivilegeCfgParse(privilegecfg.PrivilegeCfgParseRead),
		privilegecfg.RolePrivilegeCfgParse(privilegecfg.RolePrivilegeCfgParseRead),
		privilegecfg.StaffRoleCfgParse(privilegecfg.StaffRoleCfgParseRead),
		repository.PivilegeRoleRepo,
		repository.PivilegeStaffRepo,
	)
	logger.Infof("====LoadPrivilegeCfg %v", res)
	return res
}

func GetStaffPrivileges(id string) []string {
	return privilege.GetStaffPrivileges(id,
		repository.PivilegeRoleRepo,
		repository.PivilegeStaffRepo)
}
