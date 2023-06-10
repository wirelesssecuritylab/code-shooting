package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"code-shooting/infra/logger"
	rs "code-shooting/infra/restserver"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"

	"code-shooting/app/service/result"
	targetapp "code-shooting/app/service/target"
	"code-shooting/domain/entity"
	repo "code-shooting/domain/repository"
	"code-shooting/domain/service/target"
	"code-shooting/infra/common"
	"code-shooting/infra/errcode"
	"code-shooting/infra/util/tools"
	"code-shooting/interface/assembler"
	"code-shooting/interface/dto"
	"time"
)

type targetPost func(action *dto.TargetAction, ctx rs.Context) error

type TargetController struct {
	BaseController
}

const (
	ActionAdd                = "add"
	ActionQuery              = "query"
	ActionRemove             = "remove"
	ActionModify             = "modify"
	ActionShoot              = "shoot"
	ActionRefresh            = "refresh"
	ActionQueryAlreadyTarget = "queryalreadytarget"
	ActionGet                = "get"
	ActionCheck              = "check"
)

func NewTargetController() *TargetController {
	return &TargetController{}
}

func (s *TargetController) Post(ctx rs.Context) error {
	body, _ := ioutil.ReadAll(ctx.Request().Body)

	var request = &dto.TargetAction{}
	err := assembler.ParseReq(body, request)
	if err != nil {
		return ctx.String(http.StatusBadRequest, fmt.Sprintf("bad request : %s", err.Error()))
	}

	tp := makeDefaultTargetPost()
	tp = makeAdd(tp)
	tp = makeQuery(tp)
	tp = makeQueryAlreadTarget(tp)
	tp = makeRemove(tp)
	tp = s.makeModifyFunc(tp)
	tp = makeCheckFunc(tp)

	return tp(request, ctx)
}

func (s *TargetController) Upload(ctx rs.Context) error {
	ownerId, fileType := ctx.Param("id"), ctx.Param("type")
	name, file, _, err := s.openFormFile(ctx, _FormFileKey)
	if err != nil {
		return ctx.JSON(errcode.ToErrRsp(err))
	}
	defer file.Close()

	if fileType != entity.TypeTarget && fileType != entity.TypeAnswer {
		return ctx.String(http.StatusBadRequest, fmt.Sprintf("request file type error : %v", fileType))
	}

	if err := target.GetTargetService().UploadTmpFile(file, fileType, name, ownerId); err != nil {
		return ctx.JSON(http.StatusInternalServerError, &common.Response{
			Code:   http.StatusInternalServerError,
			Result: "failure",
			Status: err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, &common.Response{Code: http.StatusOK, Result: "success"})
}

func makeDefaultTargetPost() targetPost {
	return func(action *dto.TargetAction, ctx rs.Context) error {
		return ctx.JSON(http.StatusBadRequest, &common.Response{
			Code:   http.StatusBadRequest,
			Result: "failure",
			Status: "not support action",
		})
	}
}

func makeAdd(tp targetPost) targetPost {
	return func(request *dto.TargetAction, ctx rs.Context) error {
		if request.Action == ActionAdd {
			if err := target.GetTargetService().CreateTarget(&request.Target); err != nil {
				return ctx.JSON(http.StatusInternalServerError, &common.Response{
					Code:   http.StatusInternalServerError,
					Result: "failure",
					Status: err.Error(),
				})
			}
			return ctx.JSON(http.StatusOK, &common.Response{Code: http.StatusOK, Result: "success"})
		}
		return tp(request, ctx)
	}
}
func makeQueryAlreadTarget(tp targetPost) targetPost {
	return func(request *dto.TargetAction, ctx rs.Context) error {
		if request.Action == ActionQueryAlreadyTarget {
			targets, err := target.GetTargetService().QueryTargets(&request.Target)
			if err != nil {
				return ctx.JSON(http.StatusInternalServerError, &common.Response{
					Code:   http.StatusInternalServerError,
					Result: "failure",
					Status: err.Error(),
				})
			}

			var alreadyTargets []entity.TargetEntity
			for temp := 0; temp < len(targets); temp++ {

				_, err := result.GetResultService().GetUserResult("0", targets[temp].Language, targets[temp].Id, request.Target.User)
				if err == nil {
					alreadyTargets = append(alreadyTargets, targets[temp])
				}

			}
			rangeRepo := repo.GetRangeRepo()
			for i := range alreadyTargets {
				convertRangIdToName(&alreadyTargets[i], rangeRepo)
			}
			resp := &common.Response{
				Result: "success",
				Detail: alreadyTargets,
			}
			return ctx.JSON(http.StatusOK, resp)

		}
		return tp(request, ctx)
	}
}

func makeQuery(tp targetPost) targetPost {
	return func(request *dto.TargetAction, ctx rs.Context) error {
		if request.Action == ActionQuery && request.Target.Name != "" {
			res, err := target.GetTargetService().QueryTargetName(&request.Target)
			if err != nil {
				return ctx.JSON(http.StatusInternalServerError, &common.Response{
					Code:   http.StatusInternalServerError,
					Result: "failure",
					Status: err.Error(),
				})
			}
			return ctx.JSON(http.StatusOK, res)
		}

		if request.Action == ActionQuery || request.Action == ActionShoot {
			targets, err := target.GetTargetService().QueryTargets(&request.Target)
			if err != nil {
				return ctx.JSON(http.StatusInternalServerError, &common.Response{
					Code:   http.StatusInternalServerError,
					Result: "failure",
					Status: err.Error(),
				})
			}

			rangeRepo := repo.GetRangeRepo()
			for i := range targets {
				convertRangIdToName(&targets[i], rangeRepo)
			}

			resp := &common.Response{
				Result: "success",
				Detail: targets,
			}
			return ctx.JSON(http.StatusOK, resp)
		}
		return tp(request, ctx)
	}
}

func makeRemove(tp targetPost) targetPost {
	return func(request *dto.TargetAction, ctx rs.Context) error {
		if request.Action == ActionRemove {
			err := target.GetTargetService().RemoveTarget(&request.Target)
			if err != nil {
				return ctx.JSON(http.StatusInternalServerError, &common.Response{
					Code:   http.StatusInternalServerError,
					Result: "failure",
					Status: err.Error(),
				})
			}
			return ctx.JSON(http.StatusOK, &common.Response{Code: http.StatusOK, Result: "success"})
		}
		return tp(request, ctx)
	}
}

func (s *TargetController) makeModifyFunc(tp targetPost) func(req *dto.TargetAction, ctx rs.Context) error {
	return func(req *dto.TargetAction, ctx rs.Context) error {
		if req.Action != ActionModify {
			return tp(req, ctx)
		}
		if req.Target.Name == "" {
			return ctx.JSON(errcode.ToErrRsp(errors.WithMessage(errcode.ErrParamError, "name cannot be empty")))
		}
		err := target.GetTargetService().ModifyTarget(&req.Target)
		if err != nil {
			return ctx.JSON(errcode.ToErrRsp(err))
		}
		return ctx.JSON(http.StatusOK, errcode.SuccMsg())
	}
}

func (s *TargetController) ExportBatchTarget(ctx rs.Context) error {
	decode := json.NewDecoder(ctx.Request().Body)
	var params map[string]string
	decode.Decode(&params)
	if params["targetIds"] != "" {
		targetIds := strings.Split(params["targetIds"][0:len(params["targetIds"])-1], ",")
		if len(targetIds) > 0 {
			err := target.GetTargetService().ExportBatchTarget(targetIds)
			if err == nil {
				defer func() {
					err := os.Remove("/tmp/targets.zip")
					if err != nil {
						logger.Errorf("remove tmp result excel file failed, err: %v", err)
					}
				}()
				return ctx.File("/tmp/targets.zip")
			} else {
				return ctx.JSON(http.StatusInternalServerError, &common.Response{
					Code:   http.StatusInternalServerError,
					Result: "failed",
					Status: "export data failed",
				})
			}
		} else {
			return ctx.JSON(http.StatusOK, &common.Response{
				Code:   http.StatusOK,
				Result: "failed",
				Status: "id is null",
			})
		}
	} else {
		return ctx.JSON(http.StatusOK, &common.Response{
			Code:   http.StatusOK,
			Result: "failed",
			Status: "id is null",
		})
	}

}
func (s *TargetController) UpdateFiles(ctx rs.Context) error {
	targetId, fileType := ctx.Param("id"), ctx.Param("type")
	if !tools.IsContain([]string{entity.TypeTarget, entity.TypeAnswer}, fileType) {
		return ctx.JSON(errcode.ToErrRsp(errors.WithMessage(errcode.ErrParamError, "unknown file type")))
	}
	name, file, size, err := s.openFormFile(ctx, _FormFileKey)
	if err != nil {
		return ctx.JSON(errcode.ToErrRsp(err))
	}
	defer file.Close()
	if size > _FormFileMaxSize {
		return errors.WithMessage(errcode.ErrParamError, "file too large")
	}
	err = target.GetTargetService().UpdateFile(targetId, fileType, name, file)
	if err != nil {
		return ctx.JSON(errcode.ToErrRsp(err))
	}

	result.GetResultService().UpdateResultsByAnswerChanged(targetId)
	return ctx.JSON(http.StatusOK, errcode.SuccMsg())
}

func (s *TargetController) GetCodeFile(ctx rs.Context) error {
	targetId := ctx.Param("id")
	body, _ := ioutil.ReadAll(ctx.Request().Body)
	var request = &dto.TargetCodeFile{}
	err := assembler.ParseReq(body, request)
	if err != nil {
		return ctx.String(http.StatusBadRequest, fmt.Sprintf("bad request : %s", err.Error()))
	}
	codeFile := request.Filename
	rc, err := target.GetTargetService().GetCodeFile(targetId, codeFile)
	if err != nil {
		return ctx.JSON(errcode.ToErrRsp(err))
	}
	defer rc.Close()
	return ctx.Stream(http.StatusOK, echo.MIMETextPlain, rc)
}

func convertRangIdToName(target *entity.TargetEntity, rangeRepo repo.RangeRepo) {
	var ranges []string
	for _, rangId := range target.RelatedRanges {
		r, err := rangeRepo.Get(rangId)
		if err != nil {
			logger.Warnf("Get range %s with error %v", rangId, err)
		}
		ranges = append(ranges, r.Name)
	}
	target.RelatedRanges = ranges
}

func (s *TargetController) GetTargetAnswers(ctx rs.Context) error {
	targetId := ctx.Param("id")
	rangeId := ctx.Param("rangeId")
	if rangeId == "" {
		return ctx.JSON(errcode.ToErrRsp(errors.WithMessage(errcode.ErrParamError, "rangeId is required")))
	}

	r, err := repo.GetRangeRepo().Get(rangeId)
	if err != nil {
		return ctx.JSON(errcode.ToErrRsp(err))
	}

	resp := dto.TargetAnswers{TargetID: targetId, Answers: []*dto.TargetAnswer{}}
	if (time.Now().Unix() >= r.DesensTime.Unix()) && (r.DesensTime.Unix() != 0) && (r.DesensTime.Unix() != -62135596800) {
		answers, err := target.GetTargetService().GetTargetAnswers(targetId)
		if err != nil {
			return ctx.JSON(errcode.ToErrRsp(err))
		}
		resp.Answers = assembler.TargetAnswer2Dto(answers)
	}
	return ctx.JSON(http.StatusOK, resp)
}

func makeCheckFunc(tp targetPost) func(req *dto.TargetAction, ctx rs.Context) error {
	return func(req *dto.TargetAction, ctx rs.Context) error {
		if req.Action != ActionCheck {
			return tp(req, ctx)
		}
		res, err := targetapp.GetTargetAppService().CheckTarget(req.Target.ID)
		if err != nil {
			return ctx.JSON(http.StatusOK, &common.Response{Result: "fail", Code: 0, Detail: err.Error()})
		}
		if len(res) > 0 {
			return ctx.JSON(http.StatusOK, &common.Response{Result: "success", Code: 1, Detail: res})
		}
		return ctx.JSON(http.StatusOK, &common.Response{Result: "success", Code: 0, Detail: res})
	}
}
