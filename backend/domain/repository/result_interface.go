package repository

import (
	"code-shooting/domain/entity"

	"go.uber.org/fx"
)

type ResultRepo interface {
	GetUserRangeResult(rangeId, language string, user *entity.UserEntity) (*entity.ShootingResult, error)
	GetUserListRangeResult(rangeId, language string, users []entity.UserEntity) ([]entity.ShootingResult, error)
	Save(result *entity.ShootingResult, user *entity.UserEntity, rangeID, language string) error
	SaveTargetResult(rangeId, language, targetId string, user *entity.UserEntity, result *entity.TargetResult) error
}

var resultRepo ResultRepo

func NewResultRepo() fx.Option {
	return fx.Options(
		fx.Populate(&resultRepo),
	)
}

func GetResultRepo() ResultRepo {
	return resultRepo
}
