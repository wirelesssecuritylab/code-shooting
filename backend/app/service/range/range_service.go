package rangesvc

import (
	"path/filepath"
	"sort"
	"time"

	"code-shooting/infra/logger"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"code-shooting/domain/entity"
	"code-shooting/domain/repository"
	"code-shooting/domain/service/project"
	"code-shooting/domain/service/result"
	"code-shooting/domain/service/target"
	"code-shooting/domain/service/user"
	"code-shooting/infra/errcode"
	sr "code-shooting/infra/shooting-result"
	"code-shooting/infra/util/tools"
	"code-shooting/interface/dto"
)

const (
	RangeTypeTest    = "test"
	RangeTypeCompete = "compete"
)

type RangeService struct{}

func GetRangeService() *RangeService {
	return &RangeService{}
}

type RangeAnswers struct {
	RangeID    string `json:"rangid"`
	Language   string `json:"language"`
	TotalNum   uint32 `json:"totalNum"`
	TotalScore uint32 `json:"totalScore"`
}

type answerInfo struct {
	targetId string
	ringNum  uint32
	score    uint32
}

func (s *RangeService) AddRange(dr *dto.Range) (string, error) {
	if err := s.checkAdd(dr); err != nil {
		return "", errors.WithMessage(errcode.ErrParamError, err.Error())
	}
	er := &entity.Range{
		Id:         uuid.NewString(),
		Name:       dr.Name,
		Type:       dr.Type,
		ProjectId:  dr.Project,
		Owner:      dr.Owner,
		Targets:    dr.GetTargetIds(),
		StartTime:  time.Unix(dr.StartTime, 0),
		EndTime:    time.Unix(dr.EndTime, 0),
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
		DesensTime: time.Unix(dr.DesensTime, 0),
	}
	if err := repository.GetRangeRepo().Add(er); err != nil {
		return "", err
	}
	if err := target.GetTargetService().ModifyRelatedRanges(er.Id, []string{}, er.Targets); err != nil {
		logger.Warnf("relate new range %v to targets %v failed: %v", er.Id, er.Targets, err.Error())
	}
	return er.Id, nil
}

func (s *RangeService) ModifyRange(dr *dto.Range) error {
	if err := s.checkModify(dr); err != nil {
		return errors.WithMessage(errcode.ErrParamError, err.Error())
	}
	rg, err := repository.GetRangeRepo().Get(dr.Id)
	if err != nil {
		return err
	}
	oldTargets := rg.Targets
	if len(dr.Name) > 0 {
		rg.Name = dr.Name
	}
	if len(dr.Type) > 0 {
		rg.Type = dr.Type
	}
	if len(dr.Project) > 0 {
		rg.ProjectId = dr.Project
	}
	if len(dr.Owner) > 0 {
		rg.Owner = dr.Owner
	}
	if dr.StartTime > 0 || dr.Type == RangeTypeTest {
		rg.StartTime = time.Unix(dr.StartTime, 0)
	}
	if dr.EndTime > 0 || dr.Type == RangeTypeTest {
		rg.EndTime = time.Unix(dr.EndTime, 0)
	}
	rg.DesensTime = time.Unix(dr.DesensTime, 0)

	rg.Targets = dr.GetTargetIds()
	rg.UpdateTime = time.Now()
	if err := repository.GetRangeRepo().Update(rg); err != nil {
		return err
	}
	if err := target.GetTargetService().ModifyRelatedRanges(dr.Id, oldTargets, dr.GetTargetIds()); err != nil {
		return err
	}
	return nil
}

func (s *RangeService) checkModify(dr *dto.Range) error {
	if len(dr.Type) > 0 && !tools.IsContain([]string{RangeTypeTest, RangeTypeCompete}, dr.Type) {
		return errors.Errorf("invalid type '%v'", dr.Type)
	}
	if len(dr.Id) == 0 {
		return errors.Errorf("invalid id '%s'", dr.Id)
	}
	return nil
}

func (s *RangeService) checkAdd(dr *dto.Range) error {
	if dr.Name == "" {
		return errors.Errorf("name cannot be empty")
	}
	if !tools.IsContain([]string{RangeTypeTest, RangeTypeCompete}, dr.Type) {
		return errors.Errorf("invalid type '%v'", dr.Type)
	}
	if dr.Type == RangeTypeCompete && dr.StartTime >= dr.EndTime {
		return errors.Errorf("invalid compete time")
	}
	if len(dr.Targets) == 0 {
		return errors.Errorf("targets cannot be empty")
	}
	return nil
}

func (s *RangeService) QueryRanges(dr *dto.Range) ([]*dto.RangeDetail, error) {
	all, err := repository.GetRangeRepo().ListAll()

	if err != nil {
		return []*dto.RangeDetail{}, err
	}
	if usr, err := user.NewUserDomainService().QueryUser(&entity.UserEntity{Id: dr.User}); err == nil && len(usr.Id) > 0 {
		all = s.filterRangeByOwner(usr, all)
	}
	sort.Slice(all, func(i, j int) bool {
		if !all[i].UpdateTime.Equal(all[j].UpdateTime) {
			return all[i].UpdateTime.After(all[j].UpdateTime)
		}
		return all[i].CreateTime.After(all[j].CreateTime)
	})

	results := []*dto.RangeDetail{}
	for _, r := range all {
		if (r.Owner == dr.Owner || dr.Owner == "") && (dr.Id == "" || dr.Id == r.Id) {
			results = append(results, s.newRange(r))
		}
	}

	return results, nil
}
func (s *RangeService) filterRangeByOwner(owner *entity.UserEntity, ranges []*entity.Range) []*entity.Range {
	projectIds := project.GetProjectService().FindByUser(owner)
	res := make([]*entity.Range, 0, len(ranges))
	for i := range ranges {
		var flag bool
		for _, proId := range projectIds {
			if len(ranges[i].ProjectId) == 0 || ranges[i].ProjectId == proId {
				flag = true
				break
			}
		}
		if flag {
			res = append(res, ranges[i])
		}
	}
	return res
}

func (s *RangeService) newRange(r *entity.Range) *dto.RangeDetail {
	ownerName := ""
	ts, org, prjName := make([]dto.TargetInRange, 0, len(r.Targets)), "", ""
	for _, t := range r.Targets {
		name, lang := "", ""
		if target, err := target.GetTargetService().FindByID(t); err == nil {
			name, lang = target.Name, target.Language
		} else {
			logger.Warnf("get target %s info failed : %s", t, err.Error())
		}
		ts = append(ts, dto.TargetInRange{TargetId: t, TargetName: name, Language: lang})
	}

	if usr, err := user.NewUserDomainService().QueryUser(&entity.UserEntity{Id: r.Owner}); err == nil {
		org = usr.Department
		ownerName = usr.Name
	}

	if prj := project.GetProjectService().FindPrjsByID(r.ProjectId); prj != nil {
		prjName = prj.Name
	}

	ranGe := dto.Range{
		Id:         r.Id,
		Name:       r.Name,
		Project:    r.ProjectId,
		Type:       r.Type,
		Owner:      r.Owner,
		Targets:    ts,
		StartTime:  r.StartTime.Unix(),
		EndTime:    r.EndTime.Unix(),
		DesensTime: r.DesensTime.Unix(),
	}
	return &dto.RangeDetail{Range: ranGe, OrgID: org, PrjName: prjName, OwnerName: ownerName}
}

func (s *RangeService) QueryRange(id string) (*entity.Range, error) {
	return repository.GetRangeRepo().Get(id)
}

func (s *RangeService) RemoveRange(dr *dto.Range) error {
	if dr.Id == "" {
		return errors.WithMessage(errcode.ErrParamError, "id is empty")
	}

	rg, err := repository.GetRangeRepo().Get(dr.Id)
	if err != nil {
		logger.Warnf("remove range get repo failed: %v", dr.Id, err.Error())
		return err
	}

	if err := target.GetTargetService().ModifyRelatedRanges(dr.Id, rg.Targets, []string{}); err != nil {
		logger.Warnf("remove range moidfy repo related failed: %v", dr.Id, err.Error())
		return err
	}

	if err := repository.GetRangeRepo().Remove(dr.Id); err != nil {
		logger.Warnf("remove range repo failed: %v", err.Error())
		return err
	}
	return nil
}

func (s *RangeService) QueryUserShootedRange(userId string) ([]*dto.RangeDetail, error) {
	rangeDetails, err := s.QueryRanges(&dto.Range{User: userId})
	if err != nil {
		return nil, err
	}

	rs := result.NewResultService()
	var result []*dto.RangeDetail
	for i, rangeDetail := range rangeDetails {
		targets := rangeDetail.GetTargetIds()
		for _, t := range targets {
			targetEntity, err := target.GetTargetService().FindTarget(t)
			if err != nil {
				continue
			}
			_, err = rs.GetUserShootingResult(&entity.UserEntity{Id: userId}, rangeDetail.Id, targetEntity.Language, targetEntity.Id)
			if err == nil {
				result = append(result, rangeDetails[i])
				break
			}
		}
	}
	return result, nil
}

func (s *RangeService) QueryRangeAnswers(rangeId, language string) (*RangeAnswers, error) {
	rg, err := repository.GetRangeRepo().Get(rangeId)
	if err != nil {
		logger.Warnf("Find range by Id failed: %v", rangeId, err.Error())
		return nil, err
	}

	RangeAnswers := &RangeAnswers{RangeID: rangeId, Language: language}
	var totalNum, totalScore uint32
	answerSlice := []answerInfo{}
	var sr *sr.ResultCalculator

	for _, targetId := range rg.Targets {
		targetEntity, err := repository.GetTargetRepo().FindTarget(targetId)
		if err != nil {
			logger.Warnf("Find target by Id failed: %v", targetId, err.Error())
			return nil, err
		}

		if targetEntity.Language != language {
			continue
		}

		answerFile := filepath.Join(targetEntity.GetAnswerFileDir(), targetEntity.Answer)
		num, score, err := sr.GetAnswerRingNumAndScore(answerFile)
		if err != nil {
			logger.Warnf("Get ring nums and score by answerFile failed: %v", answerFile, err.Error())
			return nil, err
		}

		info := answerInfo{targetId, num, score}
		answerSlice = append(answerSlice, info)
		totalNum += num
		totalScore += score
	}

	RangeAnswers.TotalNum = totalNum
	RangeAnswers.TotalScore = totalScore
	logger.Infof("Get range answers by ID and Language: %v", rangeId, language, RangeAnswers, answerSlice)
	return RangeAnswers, nil
}
