package repository

import (
	"code-shooting/domain/entity"
	"code-shooting/domain/repository"
	"code-shooting/infra/po"
	"code-shooting/infra/util/database"
	"encoding/json"

	"code-shooting/infra/logger"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.uber.org/fx"
)

type ShootingNoteRepository struct {
	ShootingNoteDB po.ShootingNoteDB
}

func NewShootingNoteRepository() fx.Option {
	return fx.Options(
		fx.Provide(newShootingNoteRepository),
	)
}

func newShootingNoteRepository() repository.ShootingNoteRepo {
	return &ShootingNoteRepository{
		ShootingNoteDB: po.ShootingNoteDB{GormDB: database.DB},
	}
}

func (s *ShootingNoteRepository) Save(e *entity.ShootingNoteEntity) error {
	if ap := s.ShootingNoteEntity2Po(e); ap != nil {
		return s.ShootingNoteDB.Save(ap)
	}
	return errors.Errorf("shootingnote save failed")
}

func (s *ShootingNoteRepository) Remove(userid, targetid, rangeid string) error {
	return s.ShootingNoteDB.Delete(userid, targetid, rangeid)
}

func (s *ShootingNoteRepository) GetBy(targetid string) ([]*entity.ShootingNoteEntity, error) {
	notes, err := s.ShootingNoteDB.FindMulti(targetid)
	if err != nil {
		return []*entity.ShootingNoteEntity{}, err
	}
	return s.ShootingNotePos2Entities(notes), nil
}

func (s *ShootingNoteRepository) Get(userid, targetid, rangeid string) (*entity.ShootingNoteEntity, error) {
	notePo, err := s.ShootingNoteDB.FindByMultiId(userid, targetid, rangeid)
	if err != nil {
		return nil, err
	}
	return s.ShootingNotePo2Entity(notePo), nil
}

func (s *ShootingNoteRepository) ShootingNotePos2Entities(aps []*po.ShootingNotePo) []*entity.ShootingNoteEntity {
	results := make([]*entity.ShootingNoteEntity, 0, len(aps))
	for _, ap := range aps {
		if e := s.ShootingNotePo2Entity(ap); e != nil {
			results = append(results, e)
		}
	}
	return results
}

func (s *ShootingNoteRepository) ShootingNoteEntity2Po(e *entity.ShootingNoteEntity) *po.ShootingNotePo {
	recordbytes, err := json.Marshal(e.Datas)
	if err != nil {
		logger.Errorf("ShootingNoteEntity2Po marshal failed, err: [%v].", err)
		return nil
	}
	return &po.ShootingNotePo{
		ID:         s.newID(),
		UserID:     e.UserID,
		TargetID:   e.TargetID,
		UserName:   e.UserName,
		Records:    recordbytes,
		RangeID:    e.RangeID,
		UpdateTime: e.UpdateTime,
	}
}

func (s *ShootingNoteRepository) ShootingNotePo2Entity(a *po.ShootingNotePo) *entity.ShootingNoteEntity {
	if a == nil {
		return nil
	}

	datas := []entity.ShootingData{}
	if err := json.Unmarshal(a.Records, &datas); err != nil && len(a.Records) > 0 {
		logger.Errorf("ShootingNotePo2Entity unmarshal failed, id: %s, records:%s, err: [%v].", a.ID, string(a.Records), err)
	}
	return &entity.ShootingNoteEntity{
		UserID:     a.UserID,
		UserName:   a.UserName,
		TargetID:   a.TargetID,
		RangeID:    a.RangeID,
		Datas:      datas,
		UpdateTime: a.UpdateTime,
	}
}

func (s *ShootingNoteRepository) newID() string {
	return uuid.NewString()
}
