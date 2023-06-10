package score

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"code-shooting/domain/entity"
	"code-shooting/infra/errcode"
	"code-shooting/infra/util"
)

const (
	_StdAnswerDir      = "std_answers"
	AnswerPaperFileExt = "xlsm"
)

var scoreService ScoreService

func SetScoreService(clc Calculator) {
	scoreService = ScoreService{
		clc: clc,
	}
}

type ScoreService struct {
	clc Calculator
}

func GetScoreService() *ScoreService {
	return &scoreService
}

func (s *ScoreService) Score(user *entity.UserEntity, stdAnswer string, tp *TargetAnswerPaper) (*entity.TargetResult, error) {
	result, err := s.clc.CalculateShootingResult(stdAnswer, tp)
	if err != nil {
		return nil, errors.WithMessagef(errcode.ErrScoreAnswerFileError, "%v", err)
	}
	return result, nil
}

func (s *ScoreService) getStdAnswer(target Range) (string, error) {
	stdAnswer := filepath.Join(util.DataDir, _StdAnswerDir, target.RangeId, target.Language, "answer."+AnswerPaperFileExt)
	if _, err := os.Stat(stdAnswer); err != nil {
		if os.IsNotExist(err) {
			return "", errors.WithMessagef(errcode.ErrStdAnswerNotExist, "range %v, language %v", target.RangeId, target.Language)
		}
		return "", errors.WithMessagef(errcode.ErrFileSystemError, "get standard answer failed: %v", err)
	}
	return stdAnswer, nil
}
