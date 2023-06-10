package repository

import (
	"sync"

	staff_agg "code-shooting/domain/privilege/staff-agg"
)

var PivilegeStaffRepo = &pivilegeStaffRepo{staffs: make(map[string]staff_agg.StaffEntity), mu: new(sync.RWMutex)}

type pivilegeStaffRepo struct {
	staffs map[string]staff_agg.StaffEntity
	mu     *sync.RWMutex
}

func (s *pivilegeStaffRepo) OneById(id string) (staff_agg.StaffEntity, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for key, staff := range s.staffs {
		if key == id {
			return staff, true
		}
	}
	return staff_agg.StaffEntity{}, false
}
func (s *pivilegeStaffRepo) All() []staff_agg.StaffEntity {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var res []staff_agg.StaffEntity
	if len(s.staffs) == 0 {
		return nil
	}
	for _, staff := range s.staffs {
		res = append(res, staff)
	}
	return res
}
func (s *pivilegeStaffRepo) Save(staff staff_agg.StaffEntity) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.staffs[staff.Id] = staff
}

func (s *pivilegeStaffRepo) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.staffs = make(map[string]staff_agg.StaffEntity)
}
