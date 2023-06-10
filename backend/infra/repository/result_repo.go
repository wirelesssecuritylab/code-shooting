package repository

import (
	"encoding/json"
	"strings"

	"code-shooting/infra/database/pg/sql"
	"code-shooting/infra/logger"

	"github.com/pkg/errors"
	"go.uber.org/fx"

	"code-shooting/domain/entity"
	"code-shooting/domain/repository"
	"code-shooting/infra/po"
	"code-shooting/infra/util/database"
)

type ResultRepoImpl struct {
	db *sql.GormDB
}

func NewResultRepoImpl() fx.Option {
	return fx.Options(
		fx.Provide(newResultRepoImpl),
	)
}

func newResultRepoImpl() repository.ResultRepo {
	return &ResultRepoImpl{db: database.DB}
}

func (s *ResultRepoImpl) Save(sres *entity.ShootingResult, user *entity.UserEntity, rangeId, language string) error {
	logger.Debugf("save user %s range result, id: %s, language: %s, result: %v", user.Id, rangeId, language, sres)
	for _, tres := range sres.Targets {
		records, _ := json.Marshal(tres.TargetDetails)
		res := &po.TbResult{
			UserId:     user.Id,
			RangeId:    rangeId,
			Language:   strings.ToLower(language),
			TargetId:   tres.TargetId,
			HitNum:     tres.HitNum,
			HitScore:   tres.HitScore,
			TotalNum:   tres.TotalNum,
			TotalScore: tres.TotalScore,
			Records:    records,
		}
		resOld, err := s.getUserTargetResultPo(rangeId, tres.TargetId, user.Id)
		if err == nil {
			s.db.Model(res).Where("id=?", resOld.Id).Delete(resOld)
			ret := s.db.Create(res)
			return ret.Error
		}
		if repository.IsNotFound(err) {
			ret := s.db.Create(res)
			if ret.Error != nil {
				return ret.Error
			}
		}
	}

	return nil
}

func (s *ResultRepoImpl) GetUserRangeResult(rangeId, language string, user *entity.UserEntity) (*entity.ShootingResult, error) {
	ress, err := s.getUserRangeResultPo(rangeId, language, user.Id)
	if err != nil {
		return nil, err
	}
	rangeReusult := entity.ShootingResult{
		UserName: user.Name,
		UserId:   user.Id,
	}
	for i := range ress {
		t := &entity.TargetResult{
			TargetId:      ress[i].TargetId,
			HitNum:        ress[i].HitNum,
			HitScore:      ress[i].HitScore,
			TotalNum:      ress[i].TotalNum,
			TotalScore:    ress[i].TotalScore,
			TargetDetails: []entity.TargetDetail{},
		}
		json.Unmarshal(ress[i].Records, &t.TargetDetails)
		rangeReusult.Targets = append(rangeReusult.Targets, *t)
	}

	return &rangeReusult, nil
}

func users2ids(users []entity.UserEntity) []string {
	var res = make([]string, 0, len(users))
	for i := range users {
		res = append(res, users[i].Id)
	}
	return res
}
func users2idMap(users []entity.UserEntity) map[string]*entity.ShootingResult {
	var res = make(map[string]*entity.ShootingResult)
	for _, user := range users {
		res[user.Id] = &entity.ShootingResult{
			UserName: user.Name,
			UserId:   user.Id,
		}
	}
	return res
}
func idMap2ShootingRes(idMap map[string]*entity.ShootingResult) []entity.ShootingResult {
	var res = make([]entity.ShootingResult, 0, len(idMap))
	for _, v := range idMap {
		if v != nil && len(v.Targets) > 0 {
			res = append(res, *v)
		}
	}
	return res
}
func (s *ResultRepoImpl) GetUserListRangeResult(rangeId, language string, users []entity.UserEntity) ([]entity.ShootingResult, error) {
	logger.Debug("------------------5.1 GetUserListRangeResult ")
	ids := users2ids(users)
	idmap := users2idMap(users)
	logger.Debug("------------------5.2 GetUserListRangeResult ")
	ress, err := s.getUserListRangeResultPo(rangeId, language, ids)
	if err != nil {
		return nil, err
	}
	logger.Debug("------------------5.3 GetUserListRangeResult ")
	for i := range ress {
		t := &entity.TargetResult{
			TargetId:      ress[i].TargetId,
			HitNum:        ress[i].HitNum,
			HitScore:      ress[i].HitScore,
			TotalNum:      ress[i].TotalNum,
			TotalScore:    ress[i].TotalScore,
			TargetDetails: []entity.TargetDetail{},
		}
		json.Unmarshal(ress[i].Records, &t.TargetDetails)
		if rangeReusult, ok := idmap[ress[i].UserId]; ok {
			rangeReusult.Targets = append(rangeReusult.Targets, *t)
		}
	}
	logger.Debug("------------------5.4 GetUserListRangeResult ")
	return idMap2ShootingRes(idmap), nil
}

func (s *ResultRepoImpl) getUserRangeResultPo(rangeId, language, userId string) ([]po.TbResult, error) {
	ress := []po.TbResult{}
	ret := s.db.Find(&ress, "range_id=? and language=? and user_id=?", rangeId, strings.ToLower(language), userId)
	if ret.Error != nil {
		return nil, ret.Error
	}
	if ret.RowsAffected == 0 {
		return nil, &repository.NotFound{Err: errors.New("shooting result not found")}
	}
	return ress, nil
}
func (s *ResultRepoImpl) getUserListRangeResultPo(rangeId, language string, userIds []string) ([]po.TbResult, error) {
	ress := []po.TbResult{}
	ret := s.db.Find(&ress, "range_id=? and language=? and user_id IN ?", rangeId, strings.ToLower(language), userIds)
	if ret.Error != nil {
		return nil, ret.Error
	}
	if ret.RowsAffected == 0 {
		return nil, &repository.NotFound{Err: errors.New("shooting result not found")}
	}
	return ress, nil
}
func (s *ResultRepoImpl) getUserTargetResultPo(rangeId, targetId, userId string) (*po.TbResult, error) {
	ress := []po.TbResult{}
	ret := s.db.Find(&ress, "range_id=? and target_id=? and user_id=?", rangeId, targetId, userId)
	if ret.Error != nil {
		return nil, ret.Error
	}
	if ret.RowsAffected == 0 {
		return nil, &repository.NotFound{Err: errors.New("shooting result not found")}
	}
	return &ress[0], nil
}

func (s *ResultRepoImpl) SaveTargetResult(rangeId, language, targetId string, user *entity.UserEntity, result *entity.TargetResult) error {
	records, _ := json.Marshal(result.TargetDetails)
	res := &po.TbResult{
		UserId:     user.Id,
		RangeId:    rangeId,
		Language:   strings.ToLower(language),
		TargetId:   targetId,
		HitNum:     result.HitNum,
		HitScore:   result.HitScore,
		TotalNum:   result.TotalNum,
		TotalScore: result.TotalScore,
		Records:    records,
	}
	resOld, err := s.getUserTargetResultPo(rangeId, targetId, user.Id)
	if err == nil {
		s.db.Model(res).Where("id=?", resOld.Id).Delete(resOld)
		ret := s.db.Create(res)
		return ret.Error
	}
	if repository.IsNotFound(err) {
		ret := s.db.Create(res)
		if ret.Error != nil {
			return ret.Error
		}
	}
	return nil
}
