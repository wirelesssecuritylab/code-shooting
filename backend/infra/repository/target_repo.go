package repository

import (
	"code-shooting/domain/entity"
	"code-shooting/domain/entity/spec"
	"code-shooting/domain/repository"
	"code-shooting/infra/assembler"
	"code-shooting/infra/po"
	"code-shooting/infra/util/database"

	"go.uber.org/fx"
)

type TargetRepository struct {
	TargetDB po.TargetDB
}

func NewTargetRepository() fx.Option {
	return fx.Options(
		fx.Provide(newTargetRepository),
	)
}

func newTargetRepository() repository.TargetInterface {
	return &TargetRepository{
		TargetDB: po.TargetDB{GormDB: database.DB},
	}
}

func (m *TargetRepository) InsertTarget(entity *entity.TargetEntity) error {
	TargetPo := assembler.TargetEntity2Po(entity)
	return m.TargetDB.InsertTarget(TargetPo)
}

func (m *TargetRepository) UpdateTarget(entity *entity.TargetEntity) error {
	TargetPo := assembler.TargetEntity2Po(entity)
	return m.TargetDB.UpdateTarget(TargetPo)
}

func (m *TargetRepository) DeleteTarget(entity *entity.TargetEntity) error {
	TargetPo := assembler.TargetEntity2Po(entity)
	return m.TargetDB.DeleteTarget(TargetPo)
}

func (m *TargetRepository) FindTargets(where, value string) ([]entity.TargetEntity, error) {
	Targets, err := m.TargetDB.FindTargets(where, value)
	return assembler.TargetPos2Entities(Targets), err
}

func (m *TargetRepository) Find(sp spec.Spec) ([]entity.TargetEntity, error) {
	var targets []po.TargetPo
	query, args := specToQuery(sp)
	ret := m.TargetDB.Where(query, args...).Find(&targets)
	return assembler.TargetPos2Entities(targets), cvtErr(ret.Error)
}

func (m *TargetRepository) FindTarget(value string) (*entity.TargetEntity, error) {
	TargetPo, err := m.TargetDB.FindTarget(value)
	return assembler.TargetPo2Entity(&TargetPo), cvtErr(err)
}

func (m *TargetRepository) IsExist(entity *entity.TargetEntity) (bool, error) {
	TargetPo := assembler.TargetEntity2Po(entity)
	return m.TargetDB.IsExist(TargetPo)
}

func (m *TargetRepository) FindAll() ([]entity.TargetEntity, error) {
	all, err := m.TargetDB.FindAll()
	return assembler.TargetPos2Entities(all), err
}
