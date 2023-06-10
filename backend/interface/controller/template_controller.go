package controller

import (
	"bytes"
	"code-shooting/app/service/template"
	"code-shooting/domain/entity"
	"code-shooting/infra/common"
	"code-shooting/infra/errcode"
	"code-shooting/infra/shooting-result/defect"
	"code-shooting/infra/util"
	"code-shooting/interface/dto"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	rs "code-shooting/infra/restserver"
	"github.com/pkg/errors"
)

const (
	ActionEnable                 = "enable"
	ActionDisable                = "disable"
	ActionDelete                 = "delete"
	ActionQueryTemplate          = "queryTemplate"
	ActionQueryTemplateOpHistory = "queryTemplateOpHistory"
	ActionDownload               = "download"
)

type TemplateController struct {
	BaseController
}

type TemplateHandler func(rs.Context, *dto.TemplateModel) error

func NewTemplateController() *TemplateController {
	return &TemplateController{}
}

func (s *TemplateController) Upload(ctx rs.Context) error {
	workspace := ctx.Param("workspace")
	name, file, _, err := s.openFormFile(ctx, _FormFileKey)
	if err != nil {
		return ctx.JSON(errcode.ToErrRsp(err))
	}
	defer file.Close()

	version, err := template.GetTemplateVerFromFileName(name)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, &common.Response{
			Code:   http.StatusBadRequest,
			Result: "failure",
			Status: err.Error(),
		})
	}

	temp, err := template.GetTemplateAppService().QueryTemplateByVersion(version, workspace)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &common.Response{
			Code:   http.StatusInternalServerError,
			Result: "failure",
			Status: err.Error(),
		})
	}

	if temp != nil {
		return ctx.JSON(http.StatusConflict, &common.Response{
			Code:   http.StatusConflict,
			Result: "failure",
			Status: "模板版本已存在，请修改版本号后重新上传",
		})
	}

	bufFile := &bytes.Buffer{}
	bufCheck := io.TeeReader(file, bufFile)

	if err := defect.CheckDefectCode(bufCheck); err != nil {
		return ctx.JSON(http.StatusConflict, &common.Response{
			Code:   http.StatusConflict,
			Result: "failure",
			Status: err.Error(),
		})
	}

	operator := ctx.FormValue("operator")
	if err := template.GetTemplateAppService().UploadTemplate(bufFile, name, operator, version, workspace); err != nil {
		return ctx.JSON(http.StatusInternalServerError, &common.Response{
			Code:   http.StatusInternalServerError,
			Result: "failure",
			Status: err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, &common.Response{Code: http.StatusOK, Result: "success"})
}

func (s *TemplateController) Post(ctx rs.Context) error {
	req := &dto.TemplateAction{}
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(errcode.ToErrRsp(errors.WithMessage(errcode.ErrParamError, err.Error())))
	}
	handlers := map[string]TemplateHandler{
		ActionEnable:                 s.enableOrDisableTemplate,
		ActionDisable:                s.enableOrDisableTemplate,
		ActionDelete:                 s.deleteTemplate,
		ActionDownload:               s.downloadTemplate,
		ActionQueryTemplate:          s.queryTemplate,
		ActionQueryTemplateOpHistory: s.queryTemplateOpHistory,
	}
	handler, ok := handlers[req.Action]
	if !ok {
		return ctx.JSON(errcode.ToErrRsp(errors.WithMessagef(errcode.ErrParamError, "unknown action '%v'", req.Action)))
	}
	return handler(ctx, &req.Params)
}

func (s *TemplateController) enableOrDisableTemplate(ctx rs.Context, tm *dto.TemplateModel) error {
	if err := template.GetTemplateAppService().EnableOrDisableTemplate(tm.TempleteId, convertParams2OpHistory(tm)); err != nil {
		return ctx.JSON(errcode.ToErrRsp(err))
	}
	return ctx.JSON(http.StatusOK, &common.Response{Code: http.StatusOK, Result: "success"})
}

func (s *TemplateController) downloadTemplate(ctx rs.Context, tm *dto.TemplateModel) error {
	templateFile := fmt.Sprintf("代码打靶落地模板-%s.xlsm", tm.CurrentVersion)
	templateDir := filepath.Join(util.ConfDir, "templates")
	templateAbsPath := filepath.Join(templateDir, tm.Worksapce, templateFile)
	_, err := os.Stat(templateAbsPath)
	if err == nil {
		return ctx.File(templateAbsPath)
	} else {
		return ctx.JSON(http.StatusInternalServerError, &common.Response{
			Code:   http.StatusInternalServerError,
			Result: "failed",
			Status: "download template failed",
		})
	}

}

func (s *TemplateController) deleteTemplate(ctx rs.Context, tm *dto.TemplateModel) error {
	if err := template.GetTemplateAppService().DeleteTemplate(tm.TempleteId, convertParams2OpHistory(tm)); err != nil {
		return ctx.JSON(errcode.ToErrRsp(err))
	}
	return ctx.JSON(http.StatusOK, &common.Response{Code: http.StatusOK, Result: "success"})
}

func (s *TemplateController) queryTemplate(ctx rs.Context, tm *dto.TemplateModel) error {
	templates, err := template.GetTemplateAppService().QueryTemplate()
	if err != nil {
		return ctx.JSON(errcode.ToErrRsp(err))
	}
	return ctx.JSON(http.StatusOK, &common.Response{Code: http.StatusOK, Result: "success", Detail: templates})
}

func (s *TemplateController) queryTemplateOpHistory(ctx rs.Context, tm *dto.TemplateModel) error {
	templateOpHistorys, err := template.GetTemplateAppService().QueryTemplateOpHistory()
	if err != nil {
		return ctx.JSON(errcode.ToErrRsp(err))
	}
	return ctx.JSON(http.StatusOK, &common.Response{Code: http.StatusOK, Result: "success", Detail: templateOpHistorys})
}

func convertParams2OpHistory(tm *dto.TemplateModel) *entity.TemplateOpHistory {
	return &entity.TemplateOpHistory{
		Action:         tm.Action,
		CurrentVersion: tm.CurrentVersion,
		NextVersion:    tm.NextVersion,
		Changlog:       tm.Changlog,
		Operator:       tm.Operator,
		Workspace:      tm.Worksapce,
	}
}
