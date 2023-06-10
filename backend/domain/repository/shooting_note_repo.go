package repository

import (
	"code-shooting/domain/entity"

	"go.uber.org/fx"
)

type ShootingNoteRepo interface {
	Save(e *entity.ShootingNoteEntity) error
	Get(userid, targetid, rangeid string) (*entity.ShootingNoteEntity, error)
	GetBy(targetid string) ([]*entity.ShootingNoteEntity, error)
	Remove(userid, targetid, rangeid string) error
}

var shootingNoteRepo ShootingNoteRepo

func NewShootingNoteRepo() fx.Option {
	return fx.Options(
		fx.Populate(&shootingNoteRepo),
	)
}

func GetShootingNoteRepo() ShootingNoteRepo {
	return shootingNoteRepo
}
