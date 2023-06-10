package controller

import (
	"code-shooting/domain/entity"
	"code-shooting/domain/service/document"
	"code-shooting/interface/assembler"
	"fmt"
	"io/ioutil"
	"net/http"

	"code-shooting/infra/restserver"
)

type DocumentController struct {
}

func NewDocumentController() *DocumentController {
	return &DocumentController{}
}

func (d *DocumentController) GetDocmentsDir(ctx restserver.Context) error {
	docType := ctx.Param("type")
	result := document.NewDocumentService().GetDocumentDir(docType)
	return ctx.JSON(http.StatusOK, result)
}

func (d *DocumentController) GetDocumentDetail(ctx restserver.Context) error {
	catalogueId := ctx.Param("catalogueId")
	var request = &entity.DocumentsDetailReq{}
	copiedBody := ctx.Request().Body
	body, _ := ioutil.ReadAll(copiedBody)
	err := assembler.ParseReq(body, request)
	if err != nil {
		return ctx.String(http.StatusBadRequest, fmt.Sprintf("bad request : %s", err.Error()))
	}
	context := document.NewDocumentService().GetDocumentDetail(catalogueId, request.FilePath)
	return ctx.String(http.StatusOK, context)
}
