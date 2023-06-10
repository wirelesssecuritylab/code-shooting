package injection

import (
	role_agg "code-shooting/domain/privilege/role-agg"
)

type RoleRepo interface {
	OneById(id string) (role_agg.RoleEntity, bool)
	All() []role_agg.RoleEntity
	Save(role role_agg.RoleEntity)
	Clear()
}
