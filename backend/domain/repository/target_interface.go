package repository

import (
	"code-shooting/domain/entity"
	"code-shooting/domain/entity/spec"

	"go.uber.org/fx"
)

type TargetInterface interface {
	InsertTarget(entity *entity.TargetEntity) error
	UpdateTarget(entity *entity.TargetEntity) error
	DeleteTarget(entity *entity.TargetEntity) error
	FindTargets(where, value string) ([]entity.TargetEntity, error)
	Find(sp spec.Spec) ([]entity.TargetEntity, error)
	FindTarget(value string) (*entity.TargetEntity, error)
	IsExist(entity *entity.TargetEntity) (bool, error)
	FindAll() ([]entity.TargetEntity, error)
}

var targetRepo TargetInterface

func NewTargetRepo() fx.Option {
	return fx.Options(
		fx.Populate(&targetRepo),
	)
}

func GetTargetRepo() TargetInterface {
	return targetRepo
}
