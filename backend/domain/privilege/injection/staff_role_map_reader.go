package injection

import "code-shooting/domain/dto"

type StaffRoleMapReader interface {
	Read(filepath string) (dto.StaffRolesDto, error)
}
