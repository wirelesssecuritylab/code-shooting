package controller

import (
	"mime/multipart"
	"path/filepath"

	"code-shooting/infra/restserver"
	"github.com/pkg/errors"

	"code-shooting/infra/errcode"
)

type BaseController struct{}

func (s *BaseController) openFormFile(c restserver.Context, name string) (string, multipart.File, int64, error) {
	f, h, err := c.Request().FormFile(name)
	if err != nil {
		return "", nil, 0, errors.WithMessagef(errcode.ErrUnknownError, "recv form file error: %v", err)
	}
	return filepath.Base(h.Filename), f, h.Size, nil
}
