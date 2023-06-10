package shootingnote

import (
	"code-shooting/domain/entity"
	"code-shooting/domain/repository"
	"time"

	"code-shooting/infra/logger"
)

type ShootingNoteService struct{}

func GetShootingNoteService() *ShootingNoteService {
	return &ShootingNoteService{}
}

func (s *ShootingNoteService) SubmitShootingDatas(userid, targetid, rangeid string, shootingDatas []entity.TargetDetail) {
	datas := make([]entity.ShootingData, 0, len(shootingDatas))
	for _, d := range shootingDatas {
		datas = append(datas, entity.ShootingData{
			FileName:     d.FileName,
			StartLineNum: d.StartLineNum,
			EndLineNum:   d.EndLineNum,
			StartColNum:  d.StartColNum,
			EndColNum:    d.EndColNum,
			DefectCode:   d.DefectCode,
			Remark:       d.Remark,
			ScoreNum:     d.TargetScore,
		})
	}

	noteEntity := entity.ShootingNoteEntity{UserID: userid, TargetID: targetid, RangeID: rangeid, Datas: datas, UpdateTime: time.Now()}
	if err := s.RemoveDraft(userid, targetid, rangeid); err != nil {
		logger.Errorf("RemoveShootingDraft，%s %s %s failed: ", userid, targetid, rangeid, err.Error())
		return
	}
	if err := s.Save(&noteEntity); err != nil {
		logger.Errorf("SaveShootingNote，%s %s %s failed: ", userid, targetid, rangeid, err.Error())
	}
}

func (s *ShootingNoteService) Save(a *entity.ShootingNoteEntity) error {
	logger.Infof("SaveShootingNote，%s %s %s datas:%d", a.UserID, a.TargetID, a.RangeID, len(a.Datas))
	return repository.GetShootingNoteRepo().Save(a)
}

func (s *ShootingNoteService) LoadByTarget(targetid string) ([]*entity.ShootingNoteEntity, error) {
	return repository.GetShootingNoteRepo().GetBy(targetid)
}

func (s *ShootingNoteService) Load(userid, targetid, rangeid string) (*entity.ShootingNoteEntity, error) {
	return repository.GetShootingNoteRepo().Get(userid, targetid, rangeid)
}

func (s *ShootingNoteService) SaveDraft(a *entity.ShootingNoteEntity) error {
	logger.Infof("SaveShootingDraft，%s %s %s datas:%d", a.UserID, a.TargetID, a.RangeID, len(a.Datas))
	return repository.GetShootingDraftRepo().Save(a)
}

func (s *ShootingNoteService) LoadDraft(userid, targetid, rangeid string) (*entity.ShootingNoteEntity, error) {
	return repository.GetShootingDraftRepo().Get(userid, targetid, rangeid)
}

func (s *ShootingNoteService) RemoveDraft(userid, targetid, rangeid string) error {
	return repository.GetShootingDraftRepo().Remove(userid, targetid, rangeid)
}
