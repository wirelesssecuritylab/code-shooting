package ec

import (
	"code-shooting/domain/entity/spec"

	"go.uber.org/fx"
)

type ECRepo interface {
	Save(*EC) error
	Get(id string) (*EC, error)
	Find(sp spec.Spec) ([]*EC, error)
	Remove(id, uploader string) error
}

var ecRepo ECRepo

func NewSetECRepo() fx.Option {
	return fx.Options(
		fx.Populate(&ecRepo),
	)
}

func GetECRepo() ECRepo {
	return ecRepo
}
