package submit

import (
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"code-shooting/domain/entity"
	"code-shooting/domain/repository"
	"code-shooting/domain/service/score"
	"code-shooting/domain/service/shootingnote"
	"code-shooting/domain/service/target"
	"code-shooting/domain/service/user"
)

type SubmitService struct {
	defectcoderFactory func(string) (score.DefectCoder, error)
}

var submitService SubmitService

func SetSubmitService(coderFactory func(string) (score.DefectCoder, error)) {
	submitService = SubmitService{
		defectcoderFactory: coderFactory,
	}
}

func GetSubmitService() *SubmitService {
	return &submitService
}

func (s *SubmitService) Submit(userId, rangeId, language string, raps map[string]score.TargetAnswerPaper) error {
	user, err := user.NewUserDomainService().QueryUser(&entity.UserEntity{Id: userId})
	if err != nil {
		return errors.WithMessagef(err, "get user %v info failed", userId)
	}

	answerFiles, err := s.getTargetsAnswerFilePath(raps)
	if err != nil {
		return errors.WithMessagef(err, "get targets answer files", answerFiles)
	}

	res := &entity.ShootingResult{
		UserName: user.Name,
		UserId:   user.Id,
	}
	for tgId := range raps {
		targetAnswerPaper := raps[tgId]
		tres, err := score.GetScoreService().Score(user, answerFiles[tgId], &targetAnswerPaper)
		if err != nil {
			return errors.WithMessagef(err, "get socre of target %s", tgId)
		}
		tres.TargetId = tgId
		res.Targets = append(res.Targets, *tres)

		shootingnote.GetShootingNoteService().SubmitShootingDatas(userId, tgId, rangeId, tres.TargetDetails)
	}
	return repository.GetResultRepo().Save(res, user, rangeId, language)
}

func (s *SubmitService) getTargetsAnswerFilePath(taps map[string]score.TargetAnswerPaper) (map[string]string, error) {
	answerFiles := make(map[string]string, len(taps))
	for targetId := range taps {
		tg, err := target.GetTargetService().FindTarget(targetId)
		if err != nil {
			return nil, errors.Wrapf(err, "find target")
		}
		answerFiles[targetId] = filepath.Join(tg.GetAnswerFileDir(), tg.Answer)
	}
	return answerFiles, nil
}

func (s *SubmitService) saveFile(data io.Reader) (string, error) {
	f, err := os.CreateTemp(os.TempDir(), "code-shooting-*."+score.AnswerPaperFileExt)
	if err != nil {
		return "", err
	}
	defer f.Close()
	_, err = io.Copy(f, data)
	return f.Name(), err
}
