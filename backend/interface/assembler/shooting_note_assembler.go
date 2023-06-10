package assembler

import (
	"code-shooting/domain/entity"
	"code-shooting/domain/service/score"
	"code-shooting/interface/dto"
	"time"
)

func ShootingNoteDto2Entity(a *dto.ShootingNoteDto, coder score.DefectCoder) *entity.ShootingNoteEntity {
	datas := make([]entity.ShootingData, 0, len(a.Targets))
	for _, t := range a.Targets {
		defectCode := coder.EncodeDefect(t.DefectClass, t.DefectSubClass, t.DefectDescribe)
		if defectCode == "" {
			continue
		}
		datas = append(datas, entity.ShootingData{
			FileName:     t.FileName,
			StartLineNum: t.StartLineNum,
			EndLineNum:   t.EndLineNum,
			StartColNum:  t.StartColNum,
			EndColNum:    t.EndColNum,
			Remark:       t.Remark,
			DefectCode:   defectCode,
			ScoreNum:     t.ScoreNum,
		})
	}
	return &entity.ShootingNoteEntity{
		UserID:     a.UserId,
		TargetID:   a.TargetId,
		RangeID:    a.RangeID,
		Datas:      datas,
		UserName:   a.UserName,
		UpdateTime: time.Now(),
	}
}

func ShootingNoteEntity2Dto(e *entity.ShootingNoteEntity, coder score.DefectCoder) *dto.ShootingNoteDto {
	datas := make([]dto.ShootingResultDto, 0, len(e.Datas))
	for _, a := range e.Datas {
		class, subclass, describe := coder.DecodeDefect(a.DefectCode)
		inner := dto.SubmitTargetResult{
			TargetId:       e.TargetID,
			FileName:       a.FileName,
			StartLineNum:   a.StartLineNum,
			EndLineNum:     a.EndLineNum,
			DefectClass:    class,
			DefectSubClass: subclass,
			DefectDescribe: describe,
			StartColNum:    a.StartColNum,
			EndColNum:      a.EndColNum,
			Remark:         a.Remark,
		}
		datas = append(datas, dto.ShootingResultDto{inner, a.ScoreNum})
	}
	return &dto.ShootingNoteDto{
		UserId:   e.UserID,
		UserName: e.UserName,
		TargetId: e.TargetID,
		RangeID:  e.RangeID,
		Targets:  datas,
	}
}

func ShootingDraftDto2Entity(a *dto.ShootingDraftDto, coder score.DefectCoder) *entity.ShootingNoteEntity {
	datas := make([]entity.ShootingData, 0, len(a.Targets))
	for _, t := range a.Targets {
		defectCode := coder.EncodeDefect(t.DefectClass, t.DefectSubClass, t.DefectDescribe)
		if defectCode == "" {
			continue
		}
		datas = append(datas, entity.ShootingData{
			FileName:     t.FileName,
			StartLineNum: t.StartLineNum,
			EndLineNum:   t.EndLineNum,
			StartColNum:  t.StartColNum,
			EndColNum:    t.EndColNum,
			Remark:       t.Remark,
			DefectCode:   defectCode,
		})
	}
	return &entity.ShootingNoteEntity{
		UserID:     a.UserId,
		TargetID:   a.TargetId,
		RangeID:    a.RangeID,
		Datas:      datas,
		UserName:   a.UserName,
		UpdateTime: time.Now(),
	}
}

func ShootingDraftEntity2Dto(e *entity.ShootingNoteEntity, coder score.DefectCoder) *dto.ShootingDraftDto {
	datas := make([]dto.SubmitTargetResult, 0, len(e.Datas))
	for _, a := range e.Datas {
		class, subclass, describe := coder.DecodeDefect(a.DefectCode)
		datas = append(datas, dto.SubmitTargetResult{
			TargetId:       e.TargetID,
			FileName:       a.FileName,
			StartLineNum:   a.StartLineNum,
			EndLineNum:     a.EndLineNum,
			DefectClass:    class,
			DefectSubClass: subclass,
			DefectDescribe: describe,
			StartColNum:    a.StartColNum,
			EndColNum:      a.EndColNum,
			Remark:         a.Remark,
		})
	}
	return &dto.ShootingDraftDto{
		UserId:   e.UserID,
		UserName: e.UserName,
		TargetId: e.TargetID,
		RangeID:  e.RangeID,
		Targets:  datas,
	}
}
