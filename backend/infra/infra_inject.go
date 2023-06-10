package infra

import (
	"code-shooting/infra/repository"

	"go.uber.org/fx"
)

func NewInfra() fx.Option {
	return fx.Options(
		repository.NewResultRepoImpl(),
		repository.NewRangeRepoImpl(),
		repository.NewTargetRepository(),
		repository.NewShootingNoteRepository(),
		repository.NewShootingDraftRepository(),
		repository.NewECRepoImpl(),
	)
}
