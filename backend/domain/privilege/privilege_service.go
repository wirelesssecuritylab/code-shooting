package privilege

import (
	"code-shooting/domain/privilege/injection"
	role_agg "code-shooting/domain/privilege/role-agg"
	staff_agg "code-shooting/domain/privilege/staff-agg"
	"code-shooting/infra/util/tools"

	"code-shooting/infra/logger"
)

const defaultRoleId = "shootuser"

var defaultPrivileges = []string{"rangeView", "scoreView", "submitAnswerPaper", "fillinAnswerPaper"}

func getRolePrivileges(id string, roleRepo injection.RoleRepo) []string {
	if role, ok := roleRepo.OneById(id); ok {
		return role.GetPrivileges()
	}
	return defaultPrivileges
}

func GetStaffPrivileges(id string, roleRepo injection.RoleRepo,
	staffRepo injection.StaffRepo) []string {
	var res = getRolePrivileges(defaultRoleId, roleRepo)
	if staff, ok := staffRepo.OneById(id); ok {
		for _, roleid := range staff.GetRoleIds() {
			privileges := getRolePrivileges(roleid, roleRepo)
			res = tools.AppendInSliceWhenNotIn(res, privileges)
		}
	}
	return res
}

func LoadPrivileges(privilegeFilePath, roleMapFilePath, staffMapFilePath string,
	privilegeReader injection.PrivilegeCfgReader,
	roleMapReader injection.RolePrivilegeMapReader,
	staffMapReader injection.StaffRoleMapReader,
	roleRepo injection.RoleRepo,
	staffRepo injection.StaffRepo) error {
	privilegeDto, err := privilegeReader.Read(privilegeFilePath)
	if err != nil {
		return err
	}
	roleMapDto, err := roleMapReader.Read(roleMapFilePath)
	if err != nil {
		return err
	}
	staffMapDto, err := staffMapReader.Read(staffMapFilePath)
	if err != nil {
		return err
	}
	roleEntitys, err := role_agg.NewRoleEntitys(roleMapDto, privilegeDto)
	if err != nil {
		logger.Warnf("LoadPrivileges Warning %v", err)
	}
	staffEntitys, err := staff_agg.NewStaffEntitys(staffMapDto, roleMapDto)
	if err != nil {
		logger.Warnf("LoadPrivileges Warning %v", err)
	}
	roleRepo.Clear()
	for _, role := range roleEntitys {
		roleRepo.Save(role)
	}
	staffRepo.Clear()
	for _, staff := range staffEntitys {
		staffRepo.Save(staff)
	}
	return nil
}
