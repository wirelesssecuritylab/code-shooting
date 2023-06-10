package assembler

import (
	"code-shooting/domain/entity"
	"code-shooting/infra/po"
)

func TargetEntity2Po(entity *entity.TargetEntity) *po.TargetPo {
	return &po.TargetPo{
		Id:               entity.Id,
		Name:             entity.Name,
		Language:         entity.Language,
		Template:         entity.Template,
		Owner:            entity.Owner,
		OwnerName:        entity.OwnerName,
		IsShared:         entity.IsShared,
		TagId:            entity.TagId,
		MainCategory:     entity.TagName.MainCategory,
		SubCategory:      entity.TagName.SubCategory,
		DefectDetail:     entity.TagName.DefectDetail,
		CustomLabelInfo:  entity.CustomLabel,
		ExtendedLabel:    entity.ExtendedLabel,
		InstituteLabel:   entity.InstituteLabel,
		Answer:           entity.Answer,
		Targets:          entity.Targets,
		RelatedRanges:    entity.RelatedRanges,
		Workspace:        entity.Workspace,
		CreateTime:       entity.CreateTime,
		UpdateTime:       entity.UpdateTime,
		TotalAnswerNum:   entity.TotalAnswerNum,
		TotalAnswerScore: entity.TotalAnswerScore,
	}
}

func TargetPo2Entity(po *po.TargetPo) *entity.TargetEntity {
	return &entity.TargetEntity{
		Id:            po.Id,
		Name:          po.Name,
		Language:      po.Language,
		Template:      po.Template,
		Owner:         po.Owner,
		OwnerName:     po.OwnerName,
		IsShared:      po.IsShared,
		TagId:         po.TagId,
		Answer:        po.Answer,
		Targets:       po.Targets,
		RelatedRanges: po.RelatedRanges,
		TagName: entity.TagNameInfo{
			MainCategory: po.MainCategory,
			SubCategory:  po.SubCategory,
			DefectDetail: po.DefectDetail,
		},
		CustomLabel:      po.CustomLabelInfo,
		ExtendedLabel:    po.ExtendedLabel,
		InstituteLabel:   po.InstituteLabel,
		Workspace:        po.Workspace,
		CreateTime:       po.CreateTime,
		UpdateTime:       po.UpdateTime,
		TotalAnswerNum:   po.TotalAnswerNum,
		TotalAnswerScore: po.TotalAnswerScore,
	}
}

func TargetPos2Entities(targets []po.TargetPo) []entity.TargetEntity {
	targetEntities := make([]entity.TargetEntity, len(targets))
	for i := 0; i < len(targets); i++ {
		podEntity := TargetPo2Entity(&targets[i])
		targetEntities[i] = *podEntity
	}
	return targetEntities
}
