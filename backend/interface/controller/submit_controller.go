package controller

import (
	"code-shooting/domain/repository"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"time"

	"code-shooting/infra/logger"
	"code-shooting/infra/restserver"
	"github.com/pkg/errors"

	"code-shooting/app/service/submit"
	"code-shooting/domain/service/score"
	"code-shooting/domain/service/shootingnote"
	"code-shooting/domain/service/target"
	"code-shooting/infra/errcode"
	"code-shooting/infra/util"
	"code-shooting/interface/assembler"
	"code-shooting/interface/dto"
)

const (
	_FormFileKey     = "file"
	_FormFileMaxSize = 10 << 20 // 10MB

	_UserIdQueryKey   = "userId"
	_RangeIdQueryKey  = "rangeId"
	_LanguageQueryKey = "language"
	_TargetQueryKey   = "targetId"
)

type SubmitController struct {
	BaseController
	defectcoderFactory func(string) (score.DefectCoder, error)
}

var submitController SubmitController

func SetSubmitController(coderFactory func(string) (score.DefectCoder, error)) {
	submitController = SubmitController{
		defectcoderFactory: coderFactory,
	}
}

func GetSubmitController() *SubmitController {
	return &submitController
}

func (s *SubmitController) Submit(c restserver.Context) error {
	targetRange, err := s.parseRange(c)
	if err != nil {
		return c.JSON(errcode.ToErrRsp(err))
	}

	var rangeResult = &dto.RangeShootingResult{}
	body, _ := ioutil.ReadAll(c.Request().Body)
	err = assembler.ParseReq(body, rangeResult)
	if err != nil {
		return c.JSON(errcode.ToErrRsp(errors.Wrapf(err, "bad request")))
	}

	workspace := ""
	templateVersion := ""
	if len(rangeResult.Targets) > 0 {
		t, err := target.GetTargetService().FindByID(rangeResult.Targets[0].TargetId)
		if err != nil {
			return c.JSON(errcode.ToErrRsp(err))
		}
		workspace = t.Workspace
		templateVersion = t.Template
	}
	if workspace == "" {
		workspace = util.DefaultWorkspace
	}

	fileName, _ := GetCurrentTemplateFileName(workspace, templateVersion)
	coder, err := s.defectcoderFactory(filepath.Join(util.TemplateDir, workspace, fileName))
	if err != nil {
		return c.JSON(errcode.ToErrRsp(errors.Wrapf(err, "new defect coder")))
	}
	logger.Debugf("Submit range: %v, reuqest: %v", targetRange, rangeResult)
	rangeAnswerPapers := s.parseEachTargetPaper(rangeResult, targetRange.RangeId, targetRange.Language, coder)
	err = submit.GetSubmitService().Submit(rangeResult.UserId, targetRange.RangeId, targetRange.Language, rangeAnswerPapers)
	if err != nil {
		return c.JSON(errcode.ToErrRsp(err))
	}
	return c.JSON(http.StatusOK, errcode.SuccMsg())
}

func (s *SubmitController) parseRange(c restserver.Context) (score.Range, error) {
	rangeId := c.QueryParam(_RangeIdQueryKey)
	if rangeId == "" {
		return score.Range{}, errors.WithMessage(errcode.ErrParamError, "range id is required")
	}
	language := c.QueryParam(_LanguageQueryKey)
	if language == "" {
		return score.Range{}, errors.WithMessage(errcode.ErrParamError, "language is required")
	}
	return score.Range{RangeId: rangeId, Language: language}, nil
}

func (s *SubmitController) parseEachTargetPaper(sr *dto.RangeShootingResult, rangeId, lng string, coder score.DefectCoder) map[string]score.TargetAnswerPaper {
	raps := make(map[string]score.TargetAnswerPaper)
	for _, t := range sr.Targets {
		defectCode := coder.EncodeDefect(t.DefectClass, t.DefectSubClass, t.DefectDescribe)
		if defectCode == "" {
			logger.Warnf("Defect code not found for %s %s %s", t.DefectClass, t.DefectSubClass, t.DefectDescribe)
			continue
		}
		tr := score.TargetAnswer{
			FileName:     t.FileName,
			StartLineNum: t.StartLineNum,
			EndLineNum:   t.EndLineNum,
			StartColNum:  t.StartColNum,
			EndColNum:    t.EndColNum,
			DefectCode:   defectCode,
			Remark:       t.Remark,
		}

		if tp, ok := raps[t.TargetId]; ok {
			tp.Answers = append(tp.Answers, tr)
			raps[t.TargetId] = tp
		} else {
			raps[t.TargetId] = score.TargetAnswerPaper{
				Range: score.Range{
					RangeId:  rangeId,
					Language: lng,
				},
				Answers: []score.TargetAnswer{tr},
			}
		}
	}
	return raps
}

func (s *SubmitController) Save(c restserver.Context) error {
	body, _ := ioutil.ReadAll(c.Request().Body)

	datadto := &dto.ShootingNoteDto{}
	err := assembler.ParseReq(body, datadto)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("bad request : %s", err.Error()))
	}

	workspace := ""
	templateVersion := ""
	if len(datadto.Targets) > 0 {
		t, err := target.GetTargetService().FindByID(datadto.Targets[0].TargetId)
		if err != nil {
			return c.JSON(errcode.ToErrRsp(err))
		}
		workspace = t.Workspace
		templateVersion = t.Template
	}
	if workspace == "" {
		workspace = util.DefaultWorkspace
	}

	fileName, _ := GetCurrentTemplateFileName(workspace, templateVersion)
	coder, err := s.defectcoderFactory(filepath.Join(util.TemplateDir, workspace, fileName))
	if err != nil {
		return c.JSON(errcode.ToErrRsp(errors.Wrapf(err, "new defect coder")))
	}

	ShootingNoteEntity := assembler.ShootingNoteDto2Entity(datadto, coder)
	if err := shootingnote.GetShootingNoteService().Save(ShootingNoteEntity); err != nil {
		return c.JSON(errcode.ToErrRsp(err))
	}
	return c.JSON(http.StatusOK, errcode.SuccMsg())
}

func (s *SubmitController) Load(c restserver.Context) error {
	userId := c.QueryParam(_UserIdQueryKey)
	if userId == "" {
		return c.JSON(errcode.ToErrRsp(errors.WithMessage(errcode.ErrParamError, "user id is required")))
	}

	targetId := c.QueryParam(_TargetQueryKey)
	if targetId == "" {
		return c.JSON(errcode.ToErrRsp(errors.WithMessage(errcode.ErrParamError, "target id is required")))
	}

	rangeId := c.QueryParam(_RangeIdQueryKey)
	if rangeId == "" {
		return c.JSON(errcode.ToErrRsp(errors.WithMessage(errcode.ErrParamError, "range id is required")))
	}

	a, err := shootingnote.GetShootingNoteService().Load(userId, targetId, rangeId)
	if err != nil {
		return c.JSON(errcode.ToErrRsp(err))
	}
	if a == nil {
		return c.JSON(http.StatusOK, &dto.ShootingNoteDto{UserId: userId, TargetId: targetId, RangeID: rangeId})
	}

	workspace := ""
	t, err := target.GetTargetService().FindByID(targetId)
	if err == nil {
		workspace = t.Workspace
	}
	if workspace == "" {
		workspace = util.DefaultWorkspace
	}

	fileName, _ := GetCurrentTemplateFileName(workspace, t.Template)
	coder, err := s.defectcoderFactory(filepath.Join(util.TemplateDir, workspace, fileName))
	if err != nil {
		return c.JSON(errcode.ToErrRsp(errors.Wrapf(err, "new defect coder")))
	}
	shootingNoteDto := assembler.ShootingNoteEntity2Dto(a, coder)
	rangetemp, _ := repository.GetRangeRepo().Get(rangeId)
	if rangeId != "0" && (time.Now().Unix() < rangetemp.EndTime.Unix()) {
		for i := 0; i < len(shootingNoteDto.Targets); i++ {
			shootingNoteDto.Targets[i].ScoreNum = -1
		}
	}
	return c.JSON(http.StatusOK, shootingNoteDto)

}

func (s *SubmitController) SaveDraft(c restserver.Context) error {
	body, _ := ioutil.ReadAll(c.Request().Body)

	datadto := &dto.ShootingDraftDto{}
	err := assembler.ParseReq(body, datadto)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("bad request : %s", err.Error()))
	}

	workspace := ""
	templateVersion := ""
	if len(datadto.Targets) > 0 {
		t, err := target.GetTargetService().FindByID(datadto.TargetId)
		if err != nil {
			return c.JSON(errcode.ToErrRsp(err))
		}
		workspace = t.Workspace
		templateVersion = t.Template
	}
	if workspace == "" {
		workspace = util.DefaultWorkspace
	}

	fileName, _ := GetCurrentTemplateFileName(workspace, templateVersion)
	coder, err := s.defectcoderFactory(filepath.Join(util.TemplateDir, workspace, fileName))
	if err != nil {
		return c.JSON(errcode.ToErrRsp(errors.Wrapf(err, "new defect coder")))
	}

	draft := assembler.ShootingDraftDto2Entity(datadto, coder)
	if err := shootingnote.GetShootingNoteService().SaveDraft(draft); err != nil {
		return c.JSON(errcode.ToErrRsp(err))
	}
	return c.JSON(http.StatusOK, errcode.SuccMsg())
}

func (s *SubmitController) LoadDraft(c restserver.Context) error {
	userId := c.QueryParam(_UserIdQueryKey)
	if userId == "" {
		return c.JSON(errcode.ToErrRsp(errors.WithMessage(errcode.ErrParamError, "user id is required")))
	}

	targetId := c.QueryParam(_TargetQueryKey)
	if targetId == "" {
		return c.JSON(errcode.ToErrRsp(errors.WithMessage(errcode.ErrParamError, "target id is required")))
	}

	rangeId := c.QueryParam(_RangeIdQueryKey)
	if rangeId == "" {
		return c.JSON(errcode.ToErrRsp(errors.WithMessage(errcode.ErrParamError, "range id is required")))
	}

	a, err := shootingnote.GetShootingNoteService().LoadDraft(userId, targetId, rangeId)
	if err != nil {
		return c.JSON(errcode.ToErrRsp(err))
	}
	if a == nil {
		return c.JSON(http.StatusOK, &dto.ShootingDraftDto{UserId: userId, TargetId: targetId, RangeID: rangeId})
	}

	workspace := ""
	t, err := target.GetTargetService().FindByID(targetId)
	if err == nil {
		workspace = t.Workspace
	}
	if workspace == "" {
		workspace = util.DefaultWorkspace
	}

	fileName, _ := GetCurrentTemplateFileName(workspace, t.Template)
	coder, err := s.defectcoderFactory(filepath.Join(util.TemplateDir, workspace, fileName))
	if err != nil {
		return c.JSON(errcode.ToErrRsp(errors.Wrapf(err, "new defect coder")))
	}
	return c.JSON(http.StatusOK, assembler.ShootingDraftEntity2Dto(a, coder))
}
