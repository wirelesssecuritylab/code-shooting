package result

import (
	"fmt"
	"math"
	"path/filepath"
	"strings"

	"code-shooting/infra/logger"

	"github.com/pkg/errors"
	"github.com/xuri/excelize/v2"

	"code-shooting/domain/entity"
	"code-shooting/domain/repository"
	"code-shooting/domain/service/result"
	"code-shooting/domain/service/score"
	"code-shooting/domain/service/shootingnote"
	"code-shooting/domain/service/target"
	"code-shooting/domain/service/user"
	sr "code-shooting/infra/shooting-result"
	"code-shooting/infra/task"
	"code-shooting/infra/util"
	"code-shooting/infra/util/tools"
	"code-shooting/interface/dto"
)

const (
	userNameAxis   = "A%d"
	userIdAxis     = "B%d"
	departmentAxis = "C%d"
	hitNumAxis     = "D%d"
	hitScoreAxis   = "E%d"
	msgAxis        = "F%d"
)

var validLanguages = []string{"go", "python", "c++", "c", "java", "cpp"}

type ResultAppService struct {
	defectcoderFactory func(workspace string, templateVersion string) (score.DefectCoder, error)
	*task.SingleTask
}

var resultService *ResultAppService

func SetResultService(coderFactory func(workspace string, templateVersion string) (score.DefectCoder, error)) {
	updateResultTask, err := task.NewSingleTask(func(targetid string, callBackFunc func()) {
		if callBackFunc != nil {
			defer callBackFunc()
		}

		logger.Info("rescore - ", targetid, " begin")
		shootingDatas, err := shootingnote.GetShootingNoteService().LoadByTarget(targetid)
		if err != nil {
			logger.Errorf("load target %s failed:%s", targetid, err)
			return
		}

		targetEntity, err := target.GetTargetService().FindByID(targetid)
		if err != nil {
			logger.Errorf("find target entity %s failed:%s", targetid, err)
			return
		}
		answerFile := filepath.Join(targetEntity.GetAnswerFileDir(), targetEntity.Answer)
		logger.Info("rescore - answerFile:", answerFile, " saved shootingDatas:", len(shootingDatas))

		for _, userdata := range shootingDatas {
			if err = rescoreShootingDatas(userdata, answerFile, targetEntity); err != nil {
				logger.Error(err.Error())
			}
			logger.Info("rescore - user:", userdata.UserID, " rescore over")
		}
		logger.Info("rescore - ", targetid, " finished")
	})
	if err != nil {
		logger.Infof("init task pool for target answer update failed:%s", err)
	}
	resultService = &ResultAppService{
		defectcoderFactory: coderFactory,
		SingleTask:         updateResultTask,
	}
}

func GetResultService() *ResultAppService {
	return resultService
}

func (r *ResultAppService) GetUserResult(rangeID, language, targetId, userId string) ([]ResultRespDto, error) {
	user, err := user.NewUserDomainService().QueryUser(&entity.UserEntity{Id: userId})
	if err != nil {
		return nil, errors.Wrapf(err, "query user %s", userId)
	}

	rs := result.NewResultService()
	shootingResult, err := rs.GetUserShootingResult(&entity.UserEntity{Id: userId}, rangeID, language, targetId)

	if err != nil {
		return nil, errors.Wrapf(err, "get user shooting result")
	}

	totalNum, totalScore, targetIds, err := r.getRangeTotalInfo(rangeID, targetId, language)
	if err != nil {
		return nil, errors.Wrap(err, "get range total defect num and score failed")

	}
	tempRes := ResultRespDto{
		RangeScore: RangeLangScoreDto{
			TotalNum:   totalNum,
			TotalScore: totalScore,
		},
	}

	splitResults := splitResultByTargetID(shootingResult)
	for _, splitResult := range splitResults {
		if len(splitResult.Targets) == 0 {
			continue
		}

		if !tools.Contains(splitResult.Targets[0].TargetId, targetIds) {
			continue
		}

		coder := r.getDefectCoder(splitResult.Targets[0].TargetId)
		temp := *r.transResultToDto(splitResult, user, coder, nil)
		tempRes.Targets = append(tempRes.Targets, temp.Targets...)
		tempRes.RangeScore.HitNum += temp.RangeScore.HitNum
		tempRes.RangeScore.HitScore += temp.RangeScore.HitScore
		tempRes.RangeScore.TotalNum += temp.RangeScore.TotalNum
		tempRes.RangeScore.TotalScore += temp.RangeScore.TotalScore
	}

	return []ResultRespDto{tempRes}, nil
}
func users2map(users []entity.UserEntity) map[string]entity.UserEntity {
	var res = make(map[string]entity.UserEntity, len(users))
	for i := range users {
		res[users[i].Id] = users[i]
	}
	return res
}

func (r *ResultAppService) GetTargetCoders(targetId string) map[string]score.DefectCoder {
	rangeObj, err := repository.GetRangeRepo().Get(targetId)
	if err != nil {
		return map[string]score.DefectCoder{}
	}
	if rangeObj == nil {
		return map[string]score.DefectCoder{}
	}
	res := make(map[string]score.DefectCoder, len(rangeObj.Targets))
	for i := range rangeObj.Targets {
		coder := r.getDefectCoder(rangeObj.Targets[i])
		res[rangeObj.Targets[i]] = coder
	}
	return res
}
func (r *ResultAppService) GetDepartmentResults(rangeId, language, targetId, departments string, defects dto.Defects) ([]ResultRespDto, error) {
	users, err := user.NewUserDomainService().QueryAll()
	if err != nil {
		return nil, errors.Wrap(err, "query all users")
	}
	logger.Debug("------------------2.0 users ")
	rangeCoderMap := r.GetTargetCoders(rangeId)

	totalNum, totalScore, _, err := r.getRangeTotalInfo(rangeId, targetId, language)
	if err != nil {
		return nil, errors.Wrap(err, "get range total defect num and score failed")
	}

	logger.Debug("------------------2 users len", len(users))
	rs := result.NewResultService()
	dusers := rs.FilterDepartmentsUsers(users, departments)
	dusersMap := users2map(dusers)
	rrs := make([]ResultRespDto, 0)
	onceNum := 800
	for len(dusers) > 0 {
		var onceDusers []entity.UserEntity
		if len(dusers) <= onceNum {
			onceDusers = dusers
			dusers = nil
		} else {
			onceDusers = dusers[:onceNum]
			dusers = dusers[onceNum:]
		}
		rsList, err := rs.GetUserListShootingResult(onceDusers, rangeId, language, targetId)
		logger.Debug("------------------3.1 GetUserListShootingResult ", len(onceDusers), " ", len(rsList))
		if err != nil && !repository.IsNotFound(err) {
			return nil, errors.Wrapf(err, "query user %v %s result", onceDusers, err.Error())
		}
		if err == nil {
			for i := range rsList {
				tempRes := ResultRespDto{
					RangeScore: RangeLangScoreDto{
						TotalNum:   totalNum,
						TotalScore: totalScore,
					},
				}
				splitResults := splitResultByTargetID(&rsList[i])
				for _, splitResult := range splitResults {
					if len(splitResult.Targets) == 0 {
						continue
					}
					if user, ok := dusersMap[rsList[i].UserId]; ok {
						coder, okcoder := rangeCoderMap[splitResult.Targets[0].TargetId]
						if !okcoder {
							logger.Info("rangeCoderMap find not ", splitResult.Targets[0].TargetId)
							continue
						}

						temp := *r.transResultToDto(splitResult, &user, coder, defects)
						tempRes.UserId = temp.UserId
						tempRes.UserName = temp.UserName
						tempRes.TeamName = temp.TeamName
						tempRes.Department = temp.Department
						tempRes.CenterName = temp.CenterName
						tempRes.Targets = append(tempRes.Targets, temp.Targets...)
						tempRes.RangeScore.HitNum += temp.RangeScore.HitNum
						tempRes.RangeScore.HitScore += temp.RangeScore.HitScore
					} else {
						logger.Info("dusersMap not exist ", rsList[i].UserId)
					}
					logger.Debug("------------------3.1.2 result ")
				}
				if tempRes.UserId != "" {
					tempRes.RangeScore.HitScoreHundredth = calcHundredth(tempRes.RangeScore.TotalScore, tempRes.RangeScore.HitScore)
					tempRes.RangeScore.HitRate = calcRate(tempRes.RangeScore.TotalNum, tempRes.RangeScore.HitNum)
					rrs = append(rrs, tempRes)
				}
			}
		}
		logger.Debug("------------------3.2 GetUserListShootingResult build ")
	}

	return rrs, nil
}

func (r *ResultAppService) getRangeTotalInfo(rangeId, targetId, language string) (int, int, []string, error) {
	if isPracticeRange := rangeId == "0"; isPracticeRange {
		return 0, 0, []string{targetId}, nil
	}
	rg, err := repository.GetRangeRepo().Get(rangeId)
	if err != nil {
		logger.Warnf("Find range by Id failed: %v", rangeId, err.Error())
		return 0, 0, nil, err
	}
	var sr *sr.ResultCalculator
	var totalNum, totalScore uint32
	for _, targetId := range rg.Targets {
		targetEntity, err := repository.GetTargetRepo().FindTarget(targetId)
		if err != nil {
			logger.Warnf("Find target by Id failed: %v", targetId, err.Error())
			return 0, 0, nil, err
		}

		if targetEntity.Language != language {
			continue
		}

		if needSetTargetTotalAnswerInfo := (targetEntity.TotalAnswerNum == 0 || targetEntity.TotalAnswerScore == 0); needSetTargetTotalAnswerInfo {
			answerFile := filepath.Join(targetEntity.GetAnswerFileDir(), targetEntity.Answer)
			num, score, err := sr.GetAnswerRingNumAndScore(answerFile)
			if err != nil {
				logger.Warnf("Get ring nums and score by answerFile failed: %v", answerFile, err.Error())
				return 0, 0, nil, err
			}
			targetEntity.TotalAnswerNum = int(num)
			targetEntity.TotalAnswerScore = int(score)
			if err := repository.GetTargetRepo().UpdateTarget(targetEntity); err != nil {
				logger.Warnf("update target %s total answer num and score by answerFile failed: %v", targetEntity.Id, answerFile, err.Error())
				return 0, 0, nil, err
			}
		}

		totalNum += uint32(targetEntity.TotalAnswerNum)
		totalScore += uint32(targetEntity.TotalAnswerScore)
	}
	return int(totalNum), int(totalScore), rg.Targets, nil
}

func (r *ResultAppService) getDefectCoder(targetId string) score.DefectCoder {
	workspace := ""
	t, err := target.GetTargetService().FindByID(targetId)
	if err == nil {
		workspace = t.Workspace
	}
	if workspace == "" {
		workspace = util.DefaultWorkspace
	}

	coder, _ := r.defectcoderFactory(workspace, t.Template)
	return coder
}
func splitResultByTargetID(mixResult *entity.ShootingResult) []*entity.ShootingResult {
	var result []*entity.ShootingResult
	m := make(map[string]*entity.ShootingResult)
	for _, targetResult := range mixResult.Targets {
		if v, ok := m[targetResult.TargetId]; ok {
			v.Targets = append(v.Targets, targetResult)
			m[targetResult.TargetId] = v
		}
		m[targetResult.TargetId] = &entity.ShootingResult{
			UserName: mixResult.UserName,
			UserId:   mixResult.UserId,
			Targets:  []entity.TargetResult{targetResult},
		}
	}

	for _, v := range m {
		result = append(result, v)
	}
	return result
}

func (r *ResultAppService) ValidateRange(rangeID, rangeLanguage string) error {
	if rangeID == "" {
		return errors.Errorf("range id is empty")
	}
	if rangeLanguage == "" {
		errors.Errorf("range language is empty")
	}

	return nil
}

func (r *ResultAppService) transResultToDto(result *entity.ShootingResult, user *entity.UserEntity, coder score.DefectCoder, defects dto.Defects) *ResultRespDto {
	resultDto := &ResultRespDto{
		UserId:     user.Id,
		UserName:   user.Name,
		Department: user.Department,
		TeamName:   user.TeamName,
		Targets:    []TargetResultDto{},
		RangeScore: RangeLangScoreDto{},
	}

	for index := range result.Targets {
		targertResultDto := TargetResultDto{
			HitNum:            result.Targets[index].HitNum,
			HitScore:          result.Targets[index].HitScore,
			TotalNum:          result.Targets[index].TotalNum,
			TotalScore:        result.Targets[index].TotalScore,
			TargetId:          result.Targets[index].TargetId,
			HitScoreHundredth: calcHundredth(result.Targets[index].TotalScore, result.Targets[index].HitScore),
			HitRate:           calcRate(result.Targets[index].TotalNum, result.Targets[index].HitNum),
		}
		for _, detail := range result.Targets[index].TargetDetails {
			dClass, dSubClass, dDescribe := coder.DecodeDefect(detail.DefectCode)
			detailDto := TargetDetailDto{
				FileName:       detail.FileName,
				StartLineNum:   detail.StartLineNum,
				EndLineNum:     detail.EndLineNum,
				StartColNum:    detail.StartColNum,
				EndColNum:      detail.EndColNum,
				TargetScore:    detail.TargetScore,
				Remark:         detail.Remark,
				DefectClass:    dClass,
				DefectSubClass: dSubClass,
				DefectDescribe: dDescribe,
			}

			targertResultDto.Details = append(targertResultDto.Details, detailDto)
		}
		resultDto.Targets = append(resultDto.Targets, targertResultDto)
		resultDto.RangeScore.TotalNum += targertResultDto.TotalNum
		resultDto.RangeScore.TotalScore += targertResultDto.TotalScore
		resultDto.RangeScore.HitNum += targertResultDto.HitNum
		resultDto.RangeScore.HitScore += targertResultDto.HitScore
	}
	return resultDto
}

func calcHundredth(total, current int) int {
	if total <= 0 || current <= 0 {
		return 0
	}
	return int(math.Abs(float64(current) / float64(total) * 100))
}

func calcRate(total, current int) string {
	if total <= 0 || current <= 0 {
		return "0%"
	}
	return fmt.Sprintf("%.2f%%", float64(current)/float64(total)*100)
}

func (r *ResultAppService) TransResultsDtosToExcel(rsDtos []ResultRespDto, verbose, filename string) error {
	if verbose == "" || strings.ToLower(verbose) == "true" {
		return r.writeVerboseResults(rsDtos, filename)
	}
	return r.writeSimpleResults(rsDtos, filename)
}

func (r *ResultAppService) writeVerboseResults(rsDtos []ResultRespDto, filename string) error {
	f := excelize.NewFile()
	defer f.Close()
	for _, r := range rsDtos {
		sheetName := fmt.Sprintf("%s%s", r.UserName, r.UserId)
		f.NewSheet(sheetName)
		f.SetCellValue(sheetName, "A1", "靶子ID")
		f.SetCellValue(sheetName, "B1", "文件名")
		f.SetCellValue(sheetName, "C1", "起始行号")
		f.SetCellValue(sheetName, "D1", "结束行号")
		f.SetCellValue(sheetName, "E1", "缺陷大类")
		f.SetCellValue(sheetName, "F1", "缺陷小类")
		f.SetCellValue(sheetName, "G1", "缺陷细项")
		f.SetCellValue(sheetName, "H1", "缺陷备注")
		f.SetCellValue(sheetName, "I1", "得分")
		f.SetCellValue(sheetName, "J1", "命中靶环数/靶标靶环数")
		f.SetCellValue(sheetName, "K1", "命中靶环总分/靶标靶环总分")
		index := 2
		for _, t := range r.Targets {
			f.SetCellValue(sheetName, fmt.Sprintf("J%d", index), fmt.Sprintf("%d/%d", t.HitNum, t.TotalNum))
			f.SetCellValue(sheetName, fmt.Sprintf("K%d", index), fmt.Sprintf("%d/%d", t.HitScore, t.TotalScore))
			for _, d := range t.Details {
				f.SetCellValue(sheetName, fmt.Sprintf("A%d", index), t.TargetId)
				f.SetCellValue(sheetName, fmt.Sprintf("B%d", index), d.FileName)
				f.SetCellValue(sheetName, fmt.Sprintf("C%d", index), d.StartLineNum)
				f.SetCellValue(sheetName, fmt.Sprintf("D%d", index), d.EndLineNum)
				f.SetCellValue(sheetName, fmt.Sprintf("E%d", index), d.DefectClass)
				f.SetCellValue(sheetName, fmt.Sprintf("F%d", index), d.DefectSubClass)
				f.SetCellValue(sheetName, fmt.Sprintf("G%d", index), d.DefectDescribe)
				f.SetCellValue(sheetName, fmt.Sprintf("H%d", index), d.Remark)
				f.SetCellValue(sheetName, fmt.Sprintf("I%d", index), d.TargetScore)
				index += 1
			}
			index++
		}
	}
	f.DeleteSheet("Sheet1")

	return f.SaveAs(filename)
}

func (r *ResultAppService) writeSimpleResults(rsDtos []ResultRespDto, filename string) error {
	f := excelize.NewFile()
	defer f.Close()

	for _, rs := range rsDtos {
		r.writeSimpleRangeResults(rs, f)
		for _, t := range rs.Targets {
			sheetName := t.TargetId
			f.NewSheet(sheetName)
			rows, err := f.GetRows(sheetName)
			if err != nil {
				return errors.Wrapf(err, "get rows of sheet")
			}
			index := len(rows) + 1
			if len(rows) == 0 {
				f.SetCellValue(sheetName, "A1", "靶标靶环数")
				f.SetCellValue(sheetName, "B1", t.TotalNum)
				f.SetCellValue(sheetName, "A2", "靶标靶环总分")
				f.SetCellValue(sheetName, "B2", t.TotalScore)
				f.SetCellValue(sheetName, fmt.Sprintf("A%d", 3), "姓名")
				f.SetCellValue(sheetName, fmt.Sprintf("B%d", 3), "工号")
				f.SetCellValue(sheetName, fmt.Sprintf("C%d", 3), "部门")
				f.SetCellValue(sheetName, fmt.Sprintf("D%d", 3), "团队")
				f.SetCellValue(sheetName, fmt.Sprintf("E%d", 3), "命中靶环数")
				f.SetCellValue(sheetName, fmt.Sprintf("F%d", 3), "命中靶环总分")
				f.SetCellValue(sheetName, fmt.Sprintf("G%d", 3), "命中靶环总分（百分制）")
				f.SetCellValue(sheetName, fmt.Sprintf("H%d", 3), "命中率")
				index = 4
			}
			f.SetCellValue(sheetName, fmt.Sprintf("A%d", index), rs.UserName)
			f.SetCellValue(sheetName, fmt.Sprintf("B%d", index), rs.UserId)
			f.SetCellValue(sheetName, fmt.Sprintf("C%d", index), rs.Department)
			f.SetCellValue(sheetName, fmt.Sprintf("D%d", index), rs.TeamName)
			f.SetCellValue(sheetName, fmt.Sprintf("E%d", index), t.HitNum)
			f.SetCellValue(sheetName, fmt.Sprintf("F%d", index), t.HitScore)
			f.SetCellValue(sheetName, fmt.Sprintf("G%d", index), t.HitScoreHundredth)
			f.SetCellValue(sheetName, fmt.Sprintf("H%d", index), t.HitRate)
		}
	}
	f.DeleteSheet("Sheet1")
	return f.SaveAs(filename)
}

func (r *ResultAppService) writeSimpleRangeResults(rs ResultRespDto, f *excelize.File) error {
	sheetName := "total"
	f.NewSheet(sheetName)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return errors.Wrapf(err, "get rows of sheet")
	}
	index := len(rows) + 1
	if len(rows) == 0 {
		f.SetCellValue(sheetName, "A1", "靶标靶环数")
		f.SetCellValue(sheetName, "B1", rs.RangeScore.TotalNum)
		f.SetCellValue(sheetName, "A2", "靶标靶环总分")
		f.SetCellValue(sheetName, "B2", rs.RangeScore.TotalScore)
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", 3), "姓名")
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", 3), "工号")
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", 3), "部门")
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", 3), "团队")
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", 3), "命中靶环数")
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", 3), "命中靶环总分")
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", 3), "命中靶环总分（百分制）")
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", 3), "命中率")
		index = 4
	}
	f.SetCellValue(sheetName, fmt.Sprintf("A%d", index), rs.UserName)
	f.SetCellValue(sheetName, fmt.Sprintf("B%d", index), rs.UserId)
	f.SetCellValue(sheetName, fmt.Sprintf("C%d", index), rs.Department)
	f.SetCellValue(sheetName, fmt.Sprintf("D%d", index), rs.TeamName)
	f.SetCellValue(sheetName, fmt.Sprintf("E%d", index), rs.RangeScore.HitNum)
	f.SetCellValue(sheetName, fmt.Sprintf("F%d", index), rs.RangeScore.HitScore)
	f.SetCellValue(sheetName, fmt.Sprintf("G%d", index), rs.RangeScore.HitScoreHundredth)
	f.SetCellValue(sheetName, fmt.Sprintf("H%d", index), rs.RangeScore.HitRate)
	return nil
}

func (r *ResultAppService) UpdateResultsByAnswerChanged(targetid string) {
	if r.SingleTask != nil {
		r.SingleTask.ShowTask()
		r.SingleTask.SubmitTask(targetid)
	}
}

func rescoreShootingDatas(storedata *entity.ShootingNoteEntity, answerfile string, target *entity.TargetEntity) error {
	result, err := rescore(storedata, answerfile)
	if err != nil {
		return err
	}

	user := entity.UserEntity{Id: storedata.UserID, Name: storedata.UserName}
	shootingnote.GetShootingNoteService().SubmitShootingDatas(storedata.UserID, storedata.TargetID, storedata.RangeID, result.TargetDetails)

	if err = repository.GetResultRepo().SaveTargetResult(storedata.RangeID, target.Language, storedata.TargetID, &user, result); err != nil {
		logger.Error("rescore - ", storedata.RangeID, " ", target.Language, " ", storedata.TargetID, " save failed, err :", err.Error())
	}
	return nil
}

func rescore(storedata *entity.ShootingNoteEntity, answerfile string) (*entity.TargetResult, error) {
	shootingDatas := make([]score.TargetAnswer, 0, len(storedata.Datas))
	for _, s := range storedata.Datas {
		shootingDatas = append(shootingDatas, score.TargetAnswer{
			FileName:     s.FileName,
			StartLineNum: s.StartLineNum,
			EndLineNum:   s.EndLineNum,
			StartColNum:  s.StartColNum,
			EndColNum:    s.EndColNum,
			DefectCode:   s.DefectCode,
			Remark:       s.Remark,
		})
	}
	shootingPaper := &score.TargetAnswerPaper{Answers: shootingDatas}
	return score.GetScoreService().Score(nil, answerfile, shootingPaper)
}
