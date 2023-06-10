package validator

import (
	echo "github.com/labstack/echo/v4"
	"github.com/mcuadros/go-defaults"
)

type ValidatorBinder struct {
	validater     *Validator
	defaultBinder *echo.DefaultBinder
}

func NewValidatorBinder() *ValidatorBinder {
	return &ValidatorBinder{validater: &Validator{}, defaultBinder: &echo.DefaultBinder{}}
}

func (b *ValidatorBinder) Bind(i interface{}, c echo.Context) (err error) {

	if err := b.defaultBinder.Bind(i, c); err != nil {
		return err
	}
	//在bind后，设置相关的默认值
	defaults.SetDefaults(i)

	return b.validater.Validate(i)
}
