package assembler

import (
	"code-shooting/domain/entity"
	"code-shooting/interface/dto"
)

func IdDto2Entity(model *dto.IdModel) *entity.UserEntity {
	return &entity.UserEntity{
		Id: model.Id,
	}
}

func QRCodeDto2Entity(model *dto.QrCodeModel) *entity.QrCodeEntity {
	return &entity.QrCodeEntity{
		QrCodeKey:   model.QrCodeKey,
		QrCodeValue: model.QrCodeValue,
	}
}
