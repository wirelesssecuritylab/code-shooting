package repository

import (
	"code-shooting/infra/database/pg/sql"

	"go.uber.org/fx"

	"code-shooting/domain/entity/ec"
	"code-shooting/domain/entity/spec"
	"code-shooting/infra/po"
	"code-shooting/infra/util/database"
)

type ECRepoImpl struct {
}

func NewECRepoImpl() fx.Option {
	return fx.Options(
		fx.Provide(newECRepoImpl),
	)
}

func newECRepoImpl() ec.ECRepo {
	return &ECRepoImpl{}
}

func (s *ECRepoImpl) Save(e *ec.EC) error {
	o, err := s.Get(e.Id)
	if err == nil {
		e.ImportTime = o.ImportTime
		return s.update(e)
	}
	return s.add(e)
}

func (s *ECRepoImpl) add(e *ec.EC) error {
	tb := s.do2po(e)
	ret := s.db().Create(tb)
	return cvtErr(ret.Error)
}

func (s *ECRepoImpl) update(e *ec.EC) error {
	tb := s.do2po(e)
	ret := s.db().Table(tb.TableName()).Where("id = ? and importer = ?", tb.Id, tb.Importer).Select("*").Updates(tb)
	return cvtErr(ret.Error)
}

func (s *ECRepoImpl) Get(ecId string) (*ec.EC, error) {
	tb := &po.TbEC{}
	ret := s.db().Where("id = ?", ecId).First(tb)
	if ret.Error != nil {
		return nil, cvtErr(ret.Error)
	}
	return s.po2do(tb), nil
}

func (s *ECRepoImpl) Find(sp spec.Spec) ([]*ec.EC, error) {
	var tbs []*po.TbEC
	query, args := specToQuery(sp)
	ret := s.db().Table((&po.TbEC{}).TableName()).Where(query, args...).Find(&tbs)
	if ret.Error != nil {
		return nil, cvtErr(ret.Error)
	}
	es := make([]*ec.EC, 0, len(tbs))
	for _, tb := range tbs {
		es = append(es, s.po2do(tb))
	}
	return es, nil
}

func (s *ECRepoImpl) Remove(id, importer string) error {
	ret := s.db().Where("id = ? and importer = ?", id, importer).Delete(&po.TbEC{})
	return cvtErr(ret.Error)
}

func (s *ECRepoImpl) do2po(e *ec.EC) *po.TbEC {
	return &po.TbEC{
		Id:                e.Id,
		Institute:         e.Institute,
		Center:            e.Center,
		Department:        e.Department,
		Team:              e.Team,
		Importer:          e.Importer,
		ImportTime:        e.ImportTime,
		ConvertedToTarget: e.ConvertedToTarget,
		AssociatedTargets: e.AssociatedTargets,
	}
}

func (s *ECRepoImpl) po2do(tb *po.TbEC) *ec.EC {
	return &ec.EC{
		Id: tb.Id,
		Organization: ec.Organization{
			Institute:  tb.Institute,
			Center:     tb.Center,
			Department: tb.Department,
			Team:       tb.Team,
		},
		Importer:          tb.Importer,
		ImportTime:        tb.ImportTime,
		ConvertedToTarget: tb.ConvertedToTarget,
		AssociatedTargets: tb.AssociatedTargets,
	}
}

func (s *ECRepoImpl) db() *sql.GormDB {
	return database.DB
}
