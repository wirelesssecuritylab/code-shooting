package repository

import (
	"code-shooting/domain/entity"

	"go.uber.org/fx"
)

type RangeRepo interface {
	Add(*entity.Range) error
	Get(id string) (*entity.Range, error)
	ListAll() ([]*entity.Range, error)
	Update(*entity.Range) error
	Remove(id string) error
}

var rangeRepo RangeRepo

func NewRangeRepo() fx.Option {
	return fx.Options(
		fx.Populate(&rangeRepo),
	)
}

func GetRangeRepo() RangeRepo {
	return rangeRepo
}
