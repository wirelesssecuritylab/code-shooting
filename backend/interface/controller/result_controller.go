package controller

import (
	"code-shooting/domain/repository"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"code-shooting/infra/logger"
	rs "code-shooting/infra/restserver"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"code-shooting/app/service/result"
)

type resultGet func(rangeId, language, targetId string, ctx rs.Context) (interface{}, error)

type ResultController struct {
	resultGet resultGet
}

var resultControllerSingleton *ResultController

func init() {
	resultGet := makeDefaultResultGet()
	resultGet = makeUserResultGet(resultGet)
	resultGet = makeDepartmentResultGet(resultGet)
	resultGet = makeRangeValidate(resultGet)

	resultControllerSingleton = &ResultController{
		resultGet: resultGet,
	}
}

func GetResultController() *ResultController {
	return resultControllerSingleton
}

func (r *ResultController) Get(ctx rs.Context) error {
	rangeId := ctx.Param("id")
	language := ctx.Param("language")
	targetId := ctx.QueryParam("targetId")
	role := ctx.QueryParam("role")
	rangetemp, _ := repository.GetRangeRepo().Get(rangeId)
	result1, err := r.resultGet(rangeId, language, targetId, ctx)
	if err != nil {
		logger.Errorf("result get failed, err: [%v].", err)
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}
	if rangeId == "0" || (time.Now().Unix() >= rangetemp.EndTime.Unix()) || role == "admin" {
		return ctx.JSON(http.StatusOK, result1)
	} else {
		res, ok := result1.([]result.ResultRespDto)

		var targets result.Targets
		var res_ []result.Targets
		if ok {
			var targetidarray []result.TargetID
			for i := 0; i < len(res); i++ {
				for j := 0; j < len(res[i].Targets); j++ {
					var targetId result.TargetID
					targetId.TargetID = res[i].Targets[j].TargetId
					targetidarray = append(targetidarray, targetId)
				}

			}
			targets.Targets = targetidarray
		}
		res_ = append(res_, targets)
		return ctx.JSON(http.StatusOK, res_)
	}

}

func (r *ResultController) GetExcel(ctx rs.Context) error {
	rangeId := ctx.Param("id")
	language := ctx.Param("language")
	targetId := ctx.QueryParam("targetId")
	results, err := r.resultGet(rangeId, language, targetId, ctx)
	if err != nil {
		logger.Errorf("result get failed, err: [%v].", err)
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}

	verbose := ctx.QueryParam("verbose")
	tmpFileName := fmt.Sprintf("/tmp/%s-%s-%s.xlsx", rangeId, language, uuid.New().String())
	rsDtos := results.([]result.ResultRespDto)
	err = result.GetResultService().TransResultsDtosToExcel(rsDtos, verbose, tmpFileName)
	if err != nil {
		logger.Errorf("result trans to excel failed, err: [%v].", err)
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}
	defer func() {
		err := os.Remove(tmpFileName)
		if err != nil {
			logger.Errorf("remove tmp result excel file failed, err: %v", err)
		}
	}()
	return ctx.File(tmpFileName)
}

func makeDefaultResultGet() resultGet {
	return func(rangeId, language, targetId string, ctx rs.Context) (interface{}, error) {
		return nil, errors.New("not support get")
	}
}

func makeUserResultGet(rg resultGet) resultGet {
	return func(rangeId, language, targetId string, ctx rs.Context) (interface{}, error) {
		userId := ctx.QueryParam("userId")
		if userId == "" {
			return rg(rangeId, language, targetId, ctx)
		}

		rsDtos, err := result.GetResultService().GetUserResult(rangeId, language, targetId, userId)
		if err != nil {
			return nil, errors.Wrapf(err, "get user %s shooting result", userId)
		}
		return modifyResultsByVerbose(rsDtos, ctx), nil
	}
}

func makeDepartmentResultGet(rg resultGet) resultGet {
	return func(rangeId, language, targetId string, ctx rs.Context) (interface{}, error) {
		department := ctx.QueryParam("department")
		if department == "" {
			return rg(rangeId, language, targetId, ctx)
		}

		defects, _ := DoGetDefects(targetId, language, true)
		rsDtos, err := result.GetResultService().GetDepartmentResults(rangeId, language, targetId, department, defects)
		if err != nil {
			return nil, errors.Wrapf(err, "get departments %s shooting results", department)
		}
		return modifyResultsByVerbose(rsDtos, ctx), nil
	}
}

func modifyResultsByVerbose(rsDtos []result.ResultRespDto, ctx rs.Context) []result.ResultRespDto {
	verbose := ctx.QueryParam("verbose")
	if verbose == "" || strings.ToLower(verbose) == "true" {
		return rsDtos
	}
	for i := range rsDtos {
		for j := range rsDtos[i].Targets {
			rsDtos[i].Targets[j].Details = []result.TargetDetailDto{}
		}
	}
	return rsDtos
}

func makeRangeValidate(rg resultGet) resultGet {
	return func(rangeId, language, targetId string, ctx rs.Context) (interface{}, error) {
		err := result.GetResultService().ValidateRange(rangeId, language)
		if err != nil {
			return nil, errors.Wrapf(err, "invalid range, id: %s, language: %s", rangeId, language)
		}
		return rg(rangeId, language, targetId, ctx)
	}
}
