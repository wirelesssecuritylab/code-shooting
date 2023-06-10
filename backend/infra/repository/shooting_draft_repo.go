package repository

import (
	"code-shooting/domain/entity"
	"code-shooting/domain/repository"
	"code-shooting/infra/po"
	"code-shooting/infra/util/database"
	"encoding/json"

	"code-shooting/infra/logger"

	"github.com/pkg/errors"
	"go.uber.org/fx"
)

type ShootingDraftRepository struct {
	DB po.ShootingDraftDB
}

func NewShootingDraftRepository() fx.Option {
	return fx.Options(
		fx.Provide(newShootingDraftRepository),
	)
}

func newShootingDraftRepository() repository.ShootingDraftRepo {
	return &ShootingDraftRepository{
		DB: po.ShootingDraftDB{GormDB: database.DB},
	}
}

func (s *ShootingDraftRepository) Save(e *entity.ShootingNoteEntity) error {
	if ap := s.ShootingNoteEntity2DraftPo(e); ap != nil {
		return s.DB.Save(ap)
	}
	return errors.Errorf("shooting draft save failed")
}

func (s *ShootingDraftRepository) Remove(userid, targetid, rangeid string) error {
	return s.DB.Delete(s.idKey(userid, targetid, rangeid))
}

func (s *ShootingDraftRepository) Get(userid, targetid, rangeid string) (*entity.ShootingNoteEntity, error) {
	draftPo, err := s.DB.Get(s.idKey(userid, targetid, rangeid))
	if err != nil {
		return nil, err
	}
	return s.ShootingDraftPo2Entity(draftPo), nil
}

func (s *ShootingDraftRepository) ShootingDraftPos2Entities(aps []*po.ShootingDraftPo) []*entity.ShootingNoteEntity {
	results := make([]*entity.ShootingNoteEntity, 0, len(aps))
	for _, ap := range aps {
		if e := s.ShootingDraftPo2Entity(ap); e != nil {
			results = append(results, e)
		}
	}
	return results
}

func (s *ShootingDraftRepository) ShootingNoteEntity2DraftPo(e *entity.ShootingNoteEntity) *po.ShootingDraftPo {
	bytes, err := json.Marshal(e.Datas)
	if err != nil {
		logger.Errorf("ShootingNoteEntity2DraftPo marshal failed, err: [%v].", err)
		return nil
	}
	return &po.ShootingDraftPo{
		ID:         s.idKey(e.UserID, e.TargetID, e.RangeID),
		UserID:     e.UserID,
		TargetID:   e.TargetID,
		UserName:   e.UserName,
		Records:    bytes,
		UpdateTime: e.UpdateTime,
	}
}

func (s *ShootingDraftRepository) ShootingDraftPo2Entity(a *po.ShootingDraftPo) *entity.ShootingNoteEntity {
	if a == nil {
		return nil
	}

	datas := []entity.ShootingData{}
	if err := json.Unmarshal(a.Records, &datas); err != nil && len(a.Records) > 0 {
		logger.Errorf("ShootingDraftPo2Entity unmarshal failed, id: %s, records:%s, err: [%v].", a.ID, string(a.Records), err)
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

func (s *ShootingDraftRepository) idKey(userid, targetid, rangeid string) string {
	return userid + ":" + targetid + ":" + rangeid
}
