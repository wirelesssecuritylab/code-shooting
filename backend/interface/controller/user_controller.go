package controller

import (
	userservice "code-shooting/app/service/login/user-app"
	"code-shooting/domain/entity"
	"code-shooting/domain/service/user"
	"code-shooting/infra/common"
	"code-shooting/infra/errcode"
	"code-shooting/infra/handler"
	"code-shooting/interface/assembler"
	"code-shooting/interface/dto"
	"fmt"
	"io/ioutil"
	"net/http"

	"code-shooting/infra/logger"
	rs "code-shooting/infra/restserver"
	"github.com/pkg/errors"
)

type UserController struct{}

func NewUserController() *UserController {
	return &UserController{}
}

type PersonHandler func(rs.Context, *dto.Person) error

func (s *UserController) Post(ctx rs.Context) error {
	req := &dto.PersonAction{}
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(errcode.ToErrRsp(errors.WithMessage(errcode.ErrParamError, err.Error())))
	}
	handlers := map[string]PersonHandler{
		ActionModify:  s.modifyPerson,
		ActionQuery:   s.queryPerson,
		ActionRefresh: s.refreshPerson,
	}
	h, ok := handlers[req.Action]
	if !ok {
		return ctx.JSON(errcode.ToErrRsp(errors.WithMessagef(errcode.ErrParamError, "unknown action '%v'", req.Action)))
	}
	return h(ctx, &req.Params)
}
func (s *UserController) refreshPerson(ctx rs.Context, dp *dto.Person) error {
	if p, err := userservice.NewGetInfoService().RefreshPersonMessage(dp.Id); err != nil {
		return ctx.JSON(http.StatusBadRequest, handler.FailureNotFound(err.Error()))
	} else {
		return ctx.JSON(http.StatusOK, p)
	}
}

func (s *UserController) queryPerson(ctx rs.Context, dp *dto.Person) error {
	if p, err := userservice.NewGetInfoService().GetPerson(dp.Id); err != nil {
		return ctx.JSON(http.StatusBadRequest, handler.FailureNotFound(err.Error()))
	} else {
		return ctx.JSON(http.StatusOK, p)
	}
}

func (s *UserController) modifyPerson(ctx rs.Context, dp *dto.Person) error {
	res := userservice.NewGetInfoService().GetUserInfo(&entity.UserEntity{Id: dp.Id})
	if res.Result == handler.SUCCESS {
		userObj := ((*res).Detail).(*entity.UserEntity)
		modifyFlag := false
		if len(dp.TeamName) != 0 {
			userObj.TeamName = dp.TeamName
			modifyFlag = true
		}
		if len(dp.Email) != 0 {
			userObj.Email = dp.Email
			modifyFlag = true
		}
		if len(dp.CenterName) != 0 {
			userObj.CenterName = dp.CenterName
			modifyFlag = true
		}
		if len(dp.Department) != 0 {
			userObj.Department = dp.Department
			modifyFlag = true
		}
		if len(dp.Institute) != 0 {
			userObj.Institute = dp.Institute
			modifyFlag = true
		}
		if !modifyFlag {
			return ctx.JSON(http.StatusOK, nil)
		}
		if err := user.NewUserDomainService().ModifyUser(userObj); err != nil {
			return ctx.JSON(http.StatusInternalServerError, err.Error())
		}
		return ctx.JSON(http.StatusOK, dp)
	} else {
		return ctx.JSON(http.StatusInternalServerError, res)
	}
}
func (s *UserController) GetUserInfo(ctx rs.Context) error {
	logger.Debug("entry get info")
	defer PanicDefer(ctx)
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
	response = userservice.NewGetInfoService().GetUserInfo(assembler.IdDto2Entity(request))
	if response.Result == handler.SUCCESS {
		return ctx.JSON(http.StatusOK, response)
	} else {
		return ctx.JSON(http.StatusInternalServerError, response)
	}
}

func PanicDefer(c rs.Context) {
	if rcr := recover(); rcr != nil {
		logger.Fatal("panic(%v)", rcr)
		c.String(http.StatusInternalServerError, "not found by id")
	}
}
