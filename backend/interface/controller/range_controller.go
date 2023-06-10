package controller

import (
	"net/http"

	rs "code-shooting/infra/restserver"
	"github.com/pkg/errors"

	rangesvc "code-shooting/app/service/range"
	"code-shooting/app/service/result"
	"code-shooting/infra/common"
	"code-shooting/infra/errcode"
	"code-shooting/interface/dto"
)

type RangeController struct {
	BaseController
}

func NewRangeController() *RangeController {
	return &RangeController{}
}

type RangeHandler func(rs.Context, *dto.Range) error

func (s *RangeController) GetShootedRange(ctx rs.Context) error {
	userId := ctx.Param("id")
	ranges, err := rangesvc.GetRangeService().QueryUserShootedRange(userId)
	if err != nil {
		return ctx.JSON(errcode.ToErrRsp(err))
	}
	return ctx.JSON(http.StatusOK, &common.Response{Code: http.StatusOK, Result: "success", Detail: ranges})
}

func (s *RangeController) Post(ctx rs.Context) error {
	req := &dto.RangeAction{}
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(errcode.ToErrRsp(errors.WithMessage(errcode.ErrParamError, err.Error())))
	}
	handlers := map[string]RangeHandler{
		ActionAdd:    s.addRange,
		ActionModify: s.modifyRange,
		ActionQuery:  s.queryRanges,
		ActionRemove: s.removeRange,
		ActionGet:    s.getRange,
	}
	h, ok := handlers[req.Action]
	if !ok {
		return ctx.JSON(errcode.ToErrRsp(errors.WithMessagef(errcode.ErrParamError, "unknown action '%v'", req.Action)))
	}
	return h(ctx, &req.Params)
}

func (s *RangeController) addRange(ctx rs.Context, dr *dto.Range) error {
	id, err := rangesvc.GetRangeService().AddRange(dr)
	if err != nil {
		return ctx.JSON(errcode.ToErrRsp(err))
	}
	return ctx.JSON(http.StatusOK, struct {
		Id string `json:"id"`
	}{Id: id})
}

func (s *RangeController) getRange(ctx rs.Context, dr *dto.Range) error {
	range_, err := rangesvc.GetRangeService().QueryRange(dr.Id)
	if err != nil {
		return ctx.JSON(errcode.ToErrRsp(err))
	}
	return ctx.JSON(http.StatusOK, &common.Response{Code: http.StatusOK, Result: "success", Detail: range_})
}
func (s *RangeController) modifyRange(ctx rs.Context, dr *dto.Range) error {
	err := rangesvc.GetRangeService().ModifyRange(dr)
	if err != nil {
		return ctx.JSON(errcode.ToErrRsp(err))
	}
	return ctx.JSON(http.StatusOK, errcode.SuccMsg())
}

func (s *RangeController) queryRanges(ctx rs.Context, dr *dto.Range) error {
	ranges, err := rangesvc.GetRangeService().QueryRanges(dr)
	if err != nil {
		return ctx.JSON(errcode.ToErrRsp(err))
	}
	return ctx.JSON(http.StatusOK, &common.Response{Code: http.StatusOK, Result: "success", Detail: ranges})
}

func (s *RangeController) removeRange(ctx rs.Context, dr *dto.Range) error {
	err := rangesvc.GetRangeService().RemoveRange(dr)
	if err != nil {
		return ctx.JSON(errcode.ToErrRsp(err))
	}
	return ctx.JSON(http.StatusOK, errcode.SuccMsg())
}

func (s *RangeController) GetRangeAnswers(ctx rs.Context) error {
	rangeId := ctx.Param("id")
	language := ctx.Param("language")

	err := result.GetResultService().ValidateRange(rangeId, language)
	if err != nil {
		return ctx.JSON(errcode.ToErrRsp(err))
	}

	result, err := rangesvc.GetRangeService().QueryRangeAnswers(rangeId, language)
	if err != nil {
		return ctx.JSON(errcode.ToErrRsp(err))
	}
	return ctx.JSON(http.StatusOK, &common.Response{Code: http.StatusOK, Result: "success", Detail: result})
}
