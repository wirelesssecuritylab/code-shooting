package repository

import (
	"sync"

	role_agg "code-shooting/domain/privilege/role-agg"
)

var PivilegeRoleRepo = &pivilegeRoleRepo{roles: make(map[string]role_agg.RoleEntity), mu: new(sync.RWMutex)}

type pivilegeRoleRepo struct {
	roles map[string]role_agg.RoleEntity
	mu    *sync.RWMutex
}

func (s *pivilegeRoleRepo) OneById(id string) (role_agg.RoleEntity, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for key, Role := range s.roles {
		if key == id {
			return Role, true
		}
	}
	return role_agg.RoleEntity{}, false
}

func (s *pivilegeRoleRepo) All() []role_agg.RoleEntity {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var res []role_agg.RoleEntity
	if len(s.roles) == 0 {
		return nil
	}
	for _, role := range s.roles {
		res = append(res, role)
	}
	return res
}

func (s *pivilegeRoleRepo) Save(role role_agg.RoleEntity) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.roles[role.Id] = role
}

func (s *pivilegeRoleRepo) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.roles = make(map[string]role_agg.RoleEntity)
}
