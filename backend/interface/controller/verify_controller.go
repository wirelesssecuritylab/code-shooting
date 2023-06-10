package controller

import (
	userservice "code-shooting/app/service/login/user-app"
	"code-shooting/domain/entity"
	"code-shooting/infra/common"
	"code-shooting/infra/handler"
	"code-shooting/interface/assembler"
	"code-shooting/interface/dto"
	"fmt"
	"io/ioutil"
	"net/http"

	rs "code-shooting/infra/restserver"
)

type VerifyController struct{}

func NewVerifyController() *VerifyController {
	return &VerifyController{}
}

func (s *VerifyController) Verify(ctx rs.Context) error {
	var (
		request  = &dto.IdModel{}
		response *common.Response
	)
	copiedBody := ctx.Request().Body
	body, _ := ioutil.ReadAll(copiedBody)
	err := assembler.ParseReq(body, request)
	if err != nil {
		return ctx.String(http.StatusBadRequest, fmt.Sprintf("bad request : %s", err.Error()))
	}
	if request.Id == "" {
		return ctx.String(http.StatusBadRequest, fmt.Sprintf("request paramter is null : %v", request))
	}
	addUser := &entity.UserEntity{Id: request.Id}
	response = userservice.NewGetInfoService().GetUserInfoAndUpdate(addUser)
	if response.Result == handler.SUCCESS {
		return ctx.JSON(http.StatusOK, response)
	} else {
		return ctx.JSON(http.StatusInternalServerError, response)
	}
}
