package domain

import (
	"code-shooting/domain/entity/ec"
	"code-shooting/domain/repository"

	"go.uber.org/fx"
)

func NewDomain() fx.Option {
	return fx.Options(
		repository.NewResultRepo(),
		repository.NewRangeRepo(),
		repository.NewTargetRepo(),
		repository.NewShootingNoteRepo(),
		repository.NewShootingDraftRepo(),
		ec.NewSetECRepo(),
	)
}
