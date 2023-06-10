package injection

import "code-shooting/domain/dto"

type RolePrivilegeMapReader interface {
	Read(filepath string) (dto.RolePrivilegesDto, error)
}
