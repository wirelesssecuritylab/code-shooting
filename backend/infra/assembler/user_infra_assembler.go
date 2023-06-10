package assembler

import (
	"code-shooting/domain/entity"
	"code-shooting/infra/po"
)

func UserEntity2Po(entity *entity.UserEntity) *po.UserPo {
	return &po.UserPo{
		Id:         entity.Id,
		Name:       entity.Name,
		Department: entity.Department,
		OrgId:      entity.OrgId,
		Institute:  entity.Institute,
		Email:      entity.Email,
		TeamName:   entity.TeamName,
		CenterName: entity.CenterName,
	}
}

func UserPo2Entity(po *po.UserPo) *entity.UserEntity {
	return &entity.UserEntity{
		Id:         po.Id,
		Name:       po.Name,
		Department: po.Department,
		OrgId:      po.OrgId,
		Institute:  po.Institute,
		Email:      po.Email,
		TeamName:   po.TeamName,
		CenterName: po.CenterName,
	}
}

func UserPos2Entities(users []po.UserPo) []entity.UserEntity {
	userEntities := make([]entity.UserEntity, len(users))
	for i := 0; i < len(users); i++ {
		podEntity := UserPo2Entity(&users[i])
		userEntities[i] = *podEntity
	}
	return userEntities
}
