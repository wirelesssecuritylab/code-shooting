package injection

import "code-shooting/domain/dto"

type PrivilegeCfgReader interface {
	Read(filepath string) (dto.PrivilegesDto, error)
}
