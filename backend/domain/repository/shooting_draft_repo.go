package repository

import (
	"code-shooting/domain/entity"

	"go.uber.org/fx"
)

type ShootingDraftRepo interface {
	Save(e *entity.ShootingNoteEntity) error
	Get(userid, targetid, rangeid string) (*entity.ShootingNoteEntity, error)
	Remove(userid, targetid, rangeid string) error
}

var shootingDraftRepo ShootingDraftRepo

func NewShootingDraftRepo() fx.Option {
	return fx.Options(
		fx.Populate(&shootingDraftRepo),
	)
}

func GetShootingDraftRepo() ShootingDraftRepo {
	return shootingDraftRepo
}
