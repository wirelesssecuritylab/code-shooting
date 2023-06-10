package verifyservice

import (
	"code-shooting/domain/entity"
	"code-shooting/infra/common"
	"code-shooting/infra/errcode"
	"code-shooting/infra/handler"
)

type VerifySrv struct {
}

func NewVerifyService() *VerifySrv {
	return &VerifySrv{}
}

func (v *VerifySrv) VerifyUser(codeEntity *entity.QrCodeEntity) *common.Response {

	return handler.FailureNotFound(errcode.ErrRecordNotFound.Error())
}
