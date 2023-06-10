package result

import (
	"code-shooting/domain/entity"
	"code-shooting/domain/repository"
	"strings"

	"github.com/pkg/errors"
)

func NewResultService() *ResultService {
	return &ResultService{}
}

type ResultService struct{}

func (r *ResultService) GetUserShootingResult(user *entity.UserEntity, rangeID, language, targetId string) (*entity.ShootingResult, error) {
	result, err := repository.GetResultRepo().GetUserRangeResult(rangeID, language, user)
	if err != nil {
		if repository.IsNotFound(err) {
			return nil, err
		}
		return nil, errors.Wrapf(err, "get user shooting result, rangeId: %s, language: %s", rangeID, language)
	}

	if targetId == "" {
		return result, nil
	}

	for i := range result.Targets {
		target := result.Targets[i]
		if target.TargetId == targetId {
			result.Targets = []entity.TargetResult{target}
			return result, nil
		}
	}

	return nil, errors.Errorf("The result of target %s is not found in range %s language %s", targetId, rangeID, language)
}

func (r *ResultService) GetUserListShootingResult(users []entity.UserEntity, rangeID, language, targetId string) ([]entity.ShootingResult, error) {
	results, err := repository.GetResultRepo().GetUserListRangeResult(rangeID, language, users)
	if err != nil {
		if repository.IsNotFound(err) {
			return nil, err
		}
		return nil, errors.Wrapf(err, "get user shooting result, rangeId: %s, language: %s", rangeID, language)
	}

	if targetId == "" {
		return results, nil
	}
	for i := range results {
		for j := range results[i].Targets {
			target := results[i].Targets[j]
			if target.TargetId == targetId {
				results[i].Targets = []entity.TargetResult{target}
				break
			}
		}
	}
	return results, nil
}

func (r *ResultService) FilterDepartmentsUsers(users []entity.UserEntity, departments string) []entity.UserEntity {
	dusers := make([]entity.UserEntity, 0)
	if strings.ToLower(departments) == "all" {
		return users
	}
	dps := strings.Split(departments, ",")
	for _, user := range users {
		for _, dp := range dps {
			if user.Department == dp {
				dusers = append(dusers, user)
			}
		}
	}
	return dusers
}
