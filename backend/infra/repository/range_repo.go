package repository

import (
	"code-shooting/infra/database/pg/sql"

	"go.uber.org/fx"

	"code-shooting/domain/entity"
	"code-shooting/domain/repository"
	"code-shooting/infra/po"
	"code-shooting/infra/util/database"
)

type RangeRepoImpl struct {
	db *sql.GormDB
}

func NewRangeRepoImpl() fx.Option {
	return fx.Options(
		fx.Provide(newRangeRepoImpl),
	)
}

func newRangeRepoImpl() repository.RangeRepo {
	return &RangeRepoImpl{db: database.DB}
}

func (s *RangeRepoImpl) Add(er *entity.Range) error {
	tr := s.do2po(er)
	ret := s.db.Create(tr)
	return cvtErr(ret.Error)
}

func (s *RangeRepoImpl) Get(id string) (*entity.Range, error) {
	tr := &po.TbRange{}
	ret := s.db.Where("id = ?", id).First(tr)
	if ret.Error != nil {
		return nil, cvtErr(ret.Error)
	}
	return s.po2do(tr), nil
}

func (s *RangeRepoImpl) ListAll() ([]*entity.Range, error) {
	var trs []*po.TbRange
	ret := s.db.Table((&po.TbRange{}).TableName()).Find(&trs)
	if ret.Error != nil {
		return nil, cvtErr(ret.Error)
	}
	ers := make([]*entity.Range, 0, len(trs))
	for _, tr := range trs {
		ers = append(ers, s.po2do(tr))
	}
	return ers, nil
}

func (s *RangeRepoImpl) Update(er *entity.Range) error {
	tr := s.do2po(er)
	ret := s.db.Table(tr.TableName()).Where("id = ?", tr.Id).Select("*").Updates(tr)
	return cvtErr(ret.Error)
}

func (s *RangeRepoImpl) Remove(id string) error {
	ret := s.db.Where("id = ?", id).Delete(&po.TbRange{})
	return cvtErr(ret.Error)
}

func (s *RangeRepoImpl) do2po(er *entity.Range) *po.TbRange {
	return &po.TbRange{
		Id:         er.Id,
		Name:       er.Name,
		Type:       er.Type,
		ProjectId:  er.ProjectId,
		Owner:      er.Owner,
		Targets:    er.Targets,
		StartTime:  er.StartTime,
		EndTime:    er.EndTime,
		CreateTime: er.CreateTime,
		UpdateTime: er.UpdateTime,
		DesensTime: er.DesensTime,
	}
}

func (s *RangeRepoImpl) po2do(tr *po.TbRange) *entity.Range {
	return &entity.Range{
		Id:         tr.Id,
		Name:       tr.Name,
		Type:       tr.Type,
		ProjectId:  tr.ProjectId,
		Owner:      tr.Owner,
		Targets:    tr.Targets,
		StartTime:  tr.StartTime,
		EndTime:    tr.EndTime,
		CreateTime: tr.CreateTime,
		UpdateTime: tr.UpdateTime,
		DesensTime: tr.DesensTime,
	}
}
