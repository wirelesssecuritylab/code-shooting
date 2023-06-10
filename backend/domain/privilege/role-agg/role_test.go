package role_agg

import (
	"code-shooting/domain/dto"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRoleFactory(t *testing.T) {
	Convey("Givn give Privileges and role map ", t, func() {
		roleMap := dto.RolePrivilegesDto{Mappings: []dto.RolePrivilegesMapping{
			{Id: "1", Name: "admin", Mapping: dto.Privileges{PrivilegeNames: []string{"submitStandardAnswer",
				"submitTemplate",
				"createRange",
				"deleteRange",
				"viewRangeScore"}}},
		}}
		privileges := dto.PrivilegesDto{Mappings: []dto.RolePrivilegesMapping{
			{BelongedLayer: "rangeManage", Mapping: dto.Privileges{PrivilegeNames: []string{"submitStandardAnswer",
				"submitTemplate",
				"createRange",
				"deleteRange",
				"viewRangeScore"}}},
			{BelongedLayer: "codeShooting", Mapping: dto.Privileges{PrivilegeNames: []string{"rangeView",
				"scoreView",
				"submitAnswerPaper",
				"fillinAnswerPaper"}}},
		}}
		Convey("When construct the staff entitys ", func() {
			roles, err := NewRoleEntitys(roleMap, privileges)
			Convey("Then the statff is ok ", func() {
				So(err, ShouldBeNil)
				So(len(roles), ShouldEqual, 1)
				So(roles, ShouldContain, RoleEntity{Id: "1", Name: "admin", PrivilegeVos: []string{
					"submitStandardAnswer",
					"submitTemplate",
					"createRange",
					"deleteRange",
					"viewRangeScore"}})
			})
		})
	})
	Convey("Givn give role map privilege not in Privilege all", t, func() {
		roleMap := dto.RolePrivilegesDto{Mappings: []dto.RolePrivilegesMapping{
			{Id: "1", Name: "admin", Mapping: dto.Privileges{PrivilegeNames: []string{"submitStandardAnswer",
				"submitTemplate",
				"createRange",
				"deleteRange",
				"notexistPrivilege"}}},
		}}
		privileges := dto.PrivilegesDto{Mappings: []dto.RolePrivilegesMapping{
			{BelongedLayer: "rangeManage", Mapping: dto.Privileges{PrivilegeNames: []string{"submitStandardAnswer",
				"submitTemplate",
				"createRange",
				"deleteRange",
				"viewRangeScore"}}},
			{BelongedLayer: "codeShooting", Mapping: dto.Privileges{PrivilegeNames: []string{"rangeView",
				"scoreView",
				"submitAnswerPaper",
				"fillinAnswerPaper"}}},
		}}
		Convey("When construct the staff entitys ", func() {
			roles, err := NewRoleEntitys(roleMap, privileges)
			Convey("Then the statff is nil ", func() {
				So(err, ShouldNotBeNil)
				So(roles, ShouldContain, RoleEntity{Id: "1", Name: "admin", PrivilegeVos: []string{
					"submitStandardAnswer",
					"submitTemplate",
					"createRange",
					"deleteRange"}})
			})
		})
	})
}
