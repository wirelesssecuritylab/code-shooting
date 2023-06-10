package controller

import (
	"encoding/csv"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	rs "code-shooting/infra/restserver"
	"github.com/pkg/errors"

	ecsvc "code-shooting/app/service/ec"
	"code-shooting/domain/entity/ec"
	"code-shooting/infra/errcode"
	"code-shooting/interface/dto"
)

type ECController struct {
	BaseController
}

func NewECController() *ECController {
	return &ECController{}
}

func (s *ECController) ImportECs(c rs.Context) error {
	ei, err := s.parseImportParams(c)
	if err != nil {
		return c.JSON(errcode.ToErrRsp(err))
	}
	rc, err := s.recvCsvFile(c)
	if err != nil {
		return c.JSON(errcode.ToErrRsp(err))
	}
	defer rc.Close()
	err = ecsvc.GetECService().ImportFromCsv(ei, csv.NewReader(rc))
	if err != nil {
		return c.JSON(errcode.ToErrRsp(err))
	}
	return c.JSON(http.StatusOK, errcode.SuccMsg())
}

func (s *ECController) parseImportParams(c rs.Context) (*ecsvc.ECImportInfo, error) {
	institute := c.QueryParam("institute")
	center := c.QueryParam("center")
	userId := c.QueryParam("userId")
	if institute == "" || center == "" || userId == "" {
		return nil, errors.WithMessage(errcode.ErrParamError, "institute, center and userId cannot be empty")
	}
	return &ecsvc.ECImportInfo{
		Institute: institute,
		Center:    center,
		UserId:    userId,
	}, nil
}

func (s *ECController) recvCsvFile(c rs.Context) (io.ReadCloser, error) {
	name, file, _, err := s.openFormFile(c, _FormFileKey)
	if err != nil {
		return nil, err
	}
	if !strings.EqualFold(filepath.Ext(name), ".csv") {
		defer file.Close()
		return nil, errors.WithMessage(errcode.ErrParamError, "not .csv file")
	}
	return file, nil
}

func (s *ECController) QueryECs(c rs.Context) error {
	org, err := s.parseOrganization(c)
	if err != nil {
		return c.JSON(errcode.ToErrRsp(err))
	}
	tr, err := s.parseTimeRange(c)
	if err != nil {
		return c.JSON(errcode.ToErrRsp(err))
	}
	ecs, err := ecsvc.GetECService().QueryECs(org, tr)
	if err != nil {
		return c.JSON(errcode.ToErrRsp(err))
	}
	return c.JSON(http.StatusOK, &dto.QueryECRsp{ECs: ecs})
}

func (s *ECController) parseOrganization(c rs.Context) (*ec.Organization, error) {
	var org ec.Organization
	institute := c.QueryParam("institute")
	center := c.QueryParam("center")
	if institute == "" || center == "" {
		return nil, errors.WithMessage(errcode.ErrParamError, "institute and center cannot be empty")
	}
	org.Institute = institute
	org.Center = center
	if dept := c.QueryParam("department"); dept != "" {
		org.Department = dept
	}
	if team := c.QueryParam("team"); team != "" {
		org.Team = team
	}
	return &org, nil
}

func (s *ECController) parseTimeRange(c rs.Context) (*ecsvc.TimeFilter, error) {
	var tr ecsvc.TimeFilter
	if stStr := c.QueryParam("sinceTime"); stStr != "" {
		st, err := time.Parse(time.RFC3339, stStr)
		if err != nil {
			return nil, errors.WithMessagef(errcode.ErrParamError, "invalid since time: %v", err)
		}
		tr.SinceTime = st
	}
	if etStr := c.QueryParam("untilTime"); etStr != "" {
		et, err := time.Parse(time.RFC3339, etStr)
		if err != nil {
			return nil, errors.WithMessagef(errcode.ErrParamError, "invalid until time: %v", err)
		}
		tr.UntilTime = et
	}
	return &tr, nil
}

func (s *ECController) DeleteEC(c rs.Context) error {
	userId := c.QueryParam("userId")
	if userId == "" {
		return c.JSON(errcode.ToErrRsp(errors.WithMessage(errcode.ErrParamError, "userId cannot be empty")))
	}
	err := ecsvc.GetECService().DeleteEC(c.Param("id"), userId)
	if err != nil {
		return c.JSON(errcode.ToErrRsp(err))
	}
	return c.JSON(http.StatusOK, errcode.SuccMsg())
}
