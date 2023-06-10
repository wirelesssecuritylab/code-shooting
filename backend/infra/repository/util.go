package repository

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"code-shooting/infra/errcode"
)

func cvtErr(err error) error {
	switch errors.Cause(err) {
	case nil:
		return nil
	case gorm.ErrRecordNotFound:
		return errcode.ErrRecordNotFound
	default:
		return errors.WithMessage(errcode.ErrDBAccessError, err.Error())
	}
}
