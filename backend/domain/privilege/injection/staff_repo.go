package injection

import (
	staff_agg "code-shooting/domain/privilege/staff-agg"
)

type StaffRepo interface {
	OneById(id string) (staff_agg.StaffEntity, bool)
	All() []staff_agg.StaffEntity
	Save(role staff_agg.StaffEntity)
	Clear()
}
