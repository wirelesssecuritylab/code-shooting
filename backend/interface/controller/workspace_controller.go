package controller

import (
	"net/http"

	"code-shooting/app/service/workspace"

	"code-shooting/infra/restserver"
)

type WorkSpaceController struct{}

func NewWorkSpaceController() *WorkSpaceController {
	return &WorkSpaceController{}
}

func (w *WorkSpaceController) GetWorkSpaces(ctx restserver.Context) error {
	result := workspace.GetWorkSpaceAppService().GetWorkSpace()
	return ctx.JSON(http.StatusOK, result)
}
