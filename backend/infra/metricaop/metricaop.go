package metricaop

import (
	"bytes"
	"code-shooting/infra/shooting-result/defect"
	"code-shooting/infra/util"
	"code-shooting/interface/controller"
	"context"
	"io/ioutil"
	"math"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"

	"code-shooting/domain/entity"
	shootingresult "code-shooting/infra/shooting-result"
	"code-shooting/interface/assembler"
	"code-shooting/interface/dto"

	userservice "code-shooting/app/service/login/user-app"
	"code-shooting/app/service/result"
	"code-shooting/domain/service/target"

	metricData "code-shooting/infra/po/metric-data-po"
	metricRepo "code-shooting/infra/repository/metric-repo"

	"code-shooting/infra/logger"
	rs "code-shooting/infra/restserver"

	"github.com/google/uuid"
	"go.uber.org/fx"
)

const (
	ActionAdd    = "add"
	ActionQuery  = "query"
	ActionRemove = "remove"
	ActionModify = "modify"
	ActionShoot  = "shoot"
	Function     = "功能"
	Design       = "设计"
	Maintainable = "可维护"
	Security     = "安全性"
	Performance  = "性能"
	Reliability  = "可靠性"
)

// 打靶时长统计中的状态
type State uint32

const (
	Idling State = iota
	Timing
)

// 触发打靶时长统计的动作
const (
	startShoot   = "start"
	saveCopy     = "save"
	submitAnswer = "submit"
)

// 个人打靶信息
type stateOfshootingduration struct {
	FsmState       State
	ShootTimeStart time.Time
	LastEvent      string
}

// userId/tarid/state-某人，某个靶子打靶时所处的状态
var targetShootState = make(map[string]stateOfshootingduration)

// 从targetShootState获取某人，某个靶子打靶时所处的状态
func getTargetShootState(useid string, targetid string, rangeid string) stateOfshootingduration {
	stateKey := useid + targetid + rangeid
	state, ok := targetShootState[stateKey]
	if !ok {
		targetShootState[stateKey] = stateOfshootingduration{Idling, time.Now(), ""}
		return targetShootState[stateKey]
	} else {
		return state
	}

}

func setTargetShootState(useid string, targetid string, rangeid string, value stateOfshootingduration) {
	stateKey := useid + targetid + rangeid
	targetShootState[stateKey] = value

}

type MetricsCollector struct {
}

func NewExporterModule() fx.Option {
	return fx.Options(
		fx.Provide(NewMetricCollector),
	)
}

func NewMetricCollector() *MetricsCollector {
	return &MetricsCollector{}
}

func (mc *MetricsCollector) MetricsCollect(next rs.HandlerFunc) rs.HandlerFunc {
	return func(c rs.Context) error {
		// 获取URL并格式化
		url := getAccuratelyUri(c.Request().RequestURI, "v1", "?")
		// 根据URL，分别进入不同的度量
		ctx := context.Background()

		if matchRouterByUri(url, "/actions/target") {
			var request = &dto.TargetAction{}
			mc.parseRequest(c, request)
			switch request.Action {
			case ActionModify:
				user := mc.parseUser(request)
				originalTargets, _ := target.GetTargetService().QueryTargets(&request.Target)
				next(c)
				if c.Response().Status == 200 {
					mc.handleTargetModify(ctx, request, user, originalTargets[0])
				}
			case ActionAdd:
				next(c)
				if c.Response().Status == 200 {
					mc.handleTargetAdd(ctx, request)
				}
			case ActionRemove:
				user := mc.parseUser(request)
				targets, _ := target.GetTargetService().QueryTargets(&request.Target)
				next(c)
				if c.Response().Status == 200 {
					mc.handleTargetRemove(ctx, user, targets[0])
				}
			case ActionShoot:
				next(c)
				if c.Response().Status == 200 {
					mc.shootingTimeTotalMonitor(request.Target.User, request.Target.ID, request.Target.RangeID, startShoot)
				}
			default:
				next(c)
			}
			return nil
		}

		if matchRouterByUri(url, "/answers/save") {
			answerdto := &dto.ShootingNoteDto{}

			mc.parseRequest(c, answerdto)

			next(c)
			if c.Response().Status == 200 {
				mc.shootingTimeTotalMonitor(answerdto.UserId, answerdto.TargetId, answerdto.RangeID, saveCopy)
			}

			return nil
		}

		if matchRouterByUri(url, "/answers/submit") {
			var rangeResult = &dto.RangeShootingResult{}
			mc.parseRequest(c, rangeResult)
			next(c)
			if c.Response().Status == 200 {
				mc.shootingTotalMonitor(c, ctx, rangeResult)
				mc.targetShootingAccuracyMonitor(c, ctx, rangeResult)

				//统计每个靶子的打靶时长
				tagetIds := getTargetIdSet(rangeResult.Targets)
				for _, targetId := range tagetIds {
					mc.shootingTimeTotalMonitor(rangeResult.UserId, targetId, c.QueryParam("rangeId"), submitAnswer)

				}
			}
			return nil
		}

		if matchRouterByUri(url, "/actions/verify") {
			var request = &dto.IdModel{}
			mc.parseRequest(c, request)
			next(c)
			if c.Response().Status == 200 {
				metricRepo.NewRangeRequestRepository().Save(&metricData.RangeRequestPo{
					ID:          uuid.NewString(),
					UserID:      request.Id,
					RequestTime: time.Now(),
				})
			}
			return nil
		}

		next(c)

		return nil
	}
}

func (mc *MetricsCollector) parseRequest(c rs.Context, request interface{}) {
	body, _ := ioutil.ReadAll(c.Request().Body)
	c.Request().Body.Close()
	c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(body))

	err := assembler.ParseReq(body, request)
	if err != nil {
		logger.Infof("err %v", err)
	}
	logger.Infof("request %v", request)
}

func (mc *MetricsCollector) parseUser(request *dto.TargetAction) *entity.UserEntity {
	userID := &dto.IdModel{Id: request.Target.Owner}
	res := userservice.NewGetInfoService().GetUserInfo(assembler.IdDto2Entity(userID))
	return ((*res).Detail).(*entity.UserEntity)
}

func (mc *MetricsCollector) handleTargetAdd(ctx context.Context, request *dto.TargetAction) error {
	userID := &dto.IdModel{Id: request.Target.Owner}
	res := userservice.NewGetInfoService().GetUserInfo(assembler.IdDto2Entity(userID))
	user := ((*res).Detail).(*entity.UserEntity)

	targets, _ := target.GetTargetService().QueryTargets(&request.Target)
	tar := &entity.TargetEntity{Id: targets[0].Id}
	answerFile := tar.GetAnswerFileDir()
	fileName := filepath.Join(answerFile, request.Target.Answer)
	logger.Infof("answerFile %v", answerFile)
	answerShootingData, err := shootingresult.NewShootingResultCalculator().LoadShootingData(fileName)
	if err != nil {
		logger.Errorf("err %v", err)
		return err
	}
	logger.Infof("answerShootingData %v", answerShootingData)
	scoreCfg, err2 := shootingresult.NewShootingResultCalculator().LoadScoreConfig(fileName)
	if err2 != nil {
		logger.Errorf("err %v", err)
		return err
	}
	logger.Infof("scoreCfg %v", scoreCfg)
	targetTypeNum := mc.initTargetTypeNumMap()
	targetTypeScore := mc.initTargetTypeScoreMap()
	targetDefectStats := make(map[string]*metricData.TargetDefectStatPo, 10)
	saveTargetDefect := saveTargetDefectStat(&targetDefectStats)
	for _, data := range answerShootingData {
		targetTypeNum[data.DefectClass]++
		targetTypeScore[data.DefectClass] += int64(scoreCfg[data.DefectClass])
		saveTargetDefect(data.DefectCode, &metricData.TargetDefectStatPo{
			ID:        uuid.NewString(),
			TargetID:  tar.Id,
			DefectId:  data.DefectCode,
			DefectNum: 1,
		})
	}
	saveTargetDefectStatRepository(targetDefectStats)
	for targetType, nums := range targetTypeNum {
		metricRepo.NewRingNumRepository().Save(&metricData.RingNumPo{
			ID:         uuid.NewString(),
			OwnerID:    user.Id,
			TargetID:   tar.Id,
			DefectType: targetType,
			RingNum:    int(nums),
			RingScore:  int(targetTypeScore[targetType]),
			CreateTime: time.Now(),
		})
	}
	return nil
}

func (mc *MetricsCollector) handleTargetRemove(ctx context.Context, user *entity.UserEntity, originalTarget entity.TargetEntity) {
	err := metricRepo.NewTargetDefectStatRepository().Remove(originalTarget.Id)
	if err != nil {
		logger.Errorf("remove target defect statistic: %s", err.Error())
	}
}

func (mc *MetricsCollector) handleTargetModify(ctx context.Context, request *dto.TargetAction, user *entity.UserEntity, originalTarget entity.TargetEntity) error {
	tar := &entity.TargetEntity{Id: request.Target.ID}
	answerFile := tar.GetAnswerFileDir()
	fileName := filepath.Join(answerFile, request.Target.Answer)
	logger.Infof("answerFile %v", answerFile)
	answerShootingData, err := shootingresult.NewShootingResultCalculator().LoadShootingData(fileName)
	if err != nil {
		logger.Errorf("err %v", err)
		return err
	}
	logger.Infof("answerShootingData %v", answerShootingData)
	scoreCfg, err2 := shootingresult.NewShootingResultCalculator().LoadScoreConfig(fileName)
	if err2 != nil {
		logger.Errorf("err %v", err)
		return err
	}
	logger.Infof("scoreCfg %v", scoreCfg)
	targetTypeNum := mc.initTargetTypeNumMap()
	targetTypeScore := mc.initTargetTypeScoreMap()
	targetDefectStats := make(map[string]*metricData.TargetDefectStatPo, 10)
	saveTargetDefect := saveTargetDefectStat(&targetDefectStats)
	for _, data := range answerShootingData {
		targetTypeNum[data.DefectClass]++
		targetTypeScore[data.DefectClass] += int64(scoreCfg[data.DefectClass])
		saveTargetDefect(data.DefectCode, &metricData.TargetDefectStatPo{
			ID:        uuid.NewString(),
			TargetID:  tar.Id,
			DefectId:  data.DefectCode,
			DefectNum: 1,
		})
	}
	updateTargetDefectStatRepository(targetDefectStats)

	ringNumPos := make([]metricData.RingNumPo, 0, len(targetTypeNum))
	for targetType, nums := range targetTypeNum {
		ringNumPos = append(ringNumPos, metricData.RingNumPo{
			ID:         uuid.NewString(),
			OwnerID:    user.Id,
			TargetID:   tar.Id,
			DefectType: targetType,
			RingNum:    int(nums),
			RingScore:  int(targetTypeScore[targetType]),
			CreateTime: time.Now(),
		})
	}
	metricRepo.NewRingNumRepository().UpdateInBatch(&ringNumPos)

	return nil
}

func (mc *MetricsCollector) initTargetTypeNumMap() map[string]int64 {
	return map[string]int64{
		Function:     0,
		Design:       0,
		Maintainable: 0,
		Security:     0,
		Performance:  0,
		Reliability:  0,
	}
}

func (mc *MetricsCollector) initTargetTypeScoreMap() map[string]int64 {
	return map[string]int64{
		Function:     0,
		Design:       0,
		Maintainable: 0,
		Security:     0,
		Performance:  0,
		Reliability:  0,
	}
}

func (mc *MetricsCollector) shootingTotalMonitor(c rs.Context, ctx context.Context, request *dto.RangeShootingResult) {
	rangeId := c.QueryParam("rangeId")
	if rangeId == "" {
		return
	}
	language := c.QueryParam("language")
	if language == "" {
		return
	}
	metricRepo.NewShootingRecordRepository().Save(&metricData.ShootingRecordPo{
		ID:           uuid.NewString(),
		UserID:       request.UserId,
		RangeID:      rangeId,
		Language:     language,
		ShootingTime: time.Now(),
	})
}

func (mc *MetricsCollector) shootingTimeTotalMonitor(userid string, targetid string, rangeid string, evt string) {
	logger.Infof("userId :%s targetid %v ,rangeId %v", userid, targetid, rangeid)
	if userid == "" || targetid == "" {
		return
	}
	state := getTargetShootState(userid, targetid, rangeid)
	logger.Infof("evt:%v", evt)
	logger.Infof("FsmState:%v,LastEvent:%v,ShootTimeStart:%v", state.FsmState, state.LastEvent, state.ShootTimeStart)
	switch state.FsmState {
	case Idling:
		if evt == startShoot {
			logger.Infof("start shoot")
			state.FsmState = Timing
		} else {
			logger.Errorf("err evt %s in idle state", evt)
			return
		}

	case Timing:
		if evt == startShoot {
			logger.Infof("start shoot again")
		} else if evt == saveCopy || evt == submitAnswer {
			logger.Infof("%s answer", evt)
			if evt == submitAnswer {
				state.FsmState = Idling
			}

			timelen := int(math.Ceil(time.Since(state.ShootTimeStart).Minutes()))
			//先从db中查询记录，如果有，将db中的timelen  和变量标识的打靶时长相加，并更新记录，如果没有，插入timelen标识的时长的新纪录
			metricRepo.NewShootingDurationRepository().Update(&metricData.ShootingDurationPo{
				ID:       uuid.NewString(),
				UserID:   userid,
				RangeID:  rangeid,
				TargetID: targetid,
				EndTime:  time.Now(),
				Timelen:  int(timelen),
			})
		} else {
			logger.Errorf("err evt %s in Timing state", evt)
			return
		}
	default:
		logger.Errorf("invalid state :%v", state)
		return
	}
	state.ShootTimeStart = time.Now()
	state.LastEvent = evt
	setTargetShootState(userid, targetid, rangeid, state)
}

type UserInfo struct {
	Id         string `json:"id"`
	Department string `json:"department"`
	Institute  string `json:"institute"`
}

type ShootingAccuracy struct {
	UserId      string
	RangeId     string
	TargetId    string
	DefectClass string
	HitNum      int
	HitScore    int
}

func (mc *MetricsCollector) targetShootingAccuracyMonitor(c rs.Context, ctx context.Context, rangeResult *dto.RangeShootingResult) error {
	userId := rangeResult.UserId
	rangeId := c.QueryParam("rangeId")
	language := c.QueryParam("language")
	shootingAccuracys := make([]ShootingAccuracy, 0)
	//targetId 去重
	tagetIds := getTargetIdSet(rangeResult.Targets)
	for _, targetId := range tagetIds {
		userResults, err := result.GetResultService().GetUserResult(rangeId, language, targetId, userId)
		if err != nil {
			logger.Errorf("get user result failed, rangeId: %s, lange: %s,targetId: %s, userId: %s, detail: %v", rangeId, language, targetId, userId, err.Error())
			return err
		}
		logger.Infof("userResults %v", userResults)
		statisticResult(userId, rangeId, userResults, &shootingAccuracys)
	}
	logger.Infof("shootingAccuracys %v", shootingAccuracys)
	shootingAccuracyDbs := make([]metricData.ShootingAccuracyPo, len(shootingAccuracys))
	nowTime := time.Now()
	for index, shootingAccuracy := range shootingAccuracys {
		shootingAccuracyDbs[index] = metricData.ShootingAccuracyPo{
			ID:         uuid.NewString(),
			UserID:     shootingAccuracy.UserId,
			RangeID:    shootingAccuracy.RangeId,
			TargetID:   shootingAccuracy.TargetId,
			DefectType: shootingAccuracy.DefectClass,
			HitNum:     shootingAccuracy.HitNum,
			HitScore:   shootingAccuracy.HitScore,
			SubmitTime: nowTime,
		}
	}
	metricRepo.NewShootingAccuracyRepository().SaveInBatch(&shootingAccuracyDbs)

	return nil
}

func getTargetIdSet(s []dto.SubmitTargetResult) []string {
	targetIds := make([]string, 0)
	m := make(map[string]bool)
	for _, v := range s {
		targetId := v.TargetId
		if _, ok := m[targetId]; !ok {
			targetIds = append(targetIds, targetId)
			m[targetId] = true
		}
	}
	return targetIds
}

func statisticResult(userId, rangeId string, userResults []result.ResultRespDto, shootingAccuracys *[]ShootingAccuracy) {
	for urIndex := range userResults {
		targetResults := userResults[urIndex].Targets
		for trIndex := range targetResults {
			tId := targetResults[trIndex].TargetId
			details := targetResults[trIndex].Details
			hitNumMap, scoureMap := statisticHitNumAndScoure(details)
			for key, value := range scoureMap {
				shootAccuracy := &ShootingAccuracy{UserId: userId,
					RangeId:     rangeId,
					TargetId:    tId,
					DefectClass: key,
					HitNum:      hitNumMap[key],
					HitScore:    value}
				*shootingAccuracys = append(*shootingAccuracys, *shootAccuracy)
			}
		}
	}
	logger.Infof("shootingAccuracys %v", shootingAccuracys)

}

func statisticHitNumAndScoure(details []result.TargetDetailDto) (hitNumMap, scoureMap map[string]int) {
	hitNumMap = make(map[string]int)
	scoureMap = make(map[string]int)
	for detailIndex := range details {
		key := details[detailIndex].DefectClass
		if _, ok := scoureMap[key]; !ok {
			score := details[detailIndex].TargetScore
			if score > 0 {
				hitNumMap[key] = 1
			} else {
				hitNumMap[key] = 0
			}
			scoureMap[key] = score
		} else {
			score := details[detailIndex].TargetScore
			if score > 0 {
				hitNumMap[key]++
				scoureMap[key] += score
			}
		}
	}
	return
}

func getAccuratelyUri(str, start, end string) string {
	n := strings.Index(str, start)
	if n == -1 {
		n = 0
	} else {
		n += 2
	}
	str = string([]byte(str)[n:])
	m := strings.Index(str, end)
	if m == -1 {
		m = len(str)
	}
	str = string([]byte(str)[:m])
	return str
}

func matchRouterByUri(url, router string) bool {
	uris := strings.Split(url, "/")
	routers := strings.Split(router, "/")
	if len(uris) != len(routers) {
		return false
	}
	for index, value := range routers {
		if strings.Contains(value, ":") {
			continue
		}
		if value != uris[index] {
			return false
		}
	}
	return true
}

func saveTargetDefectStat(targetDefectStats *map[string]*metricData.TargetDefectStatPo) func(key string, value *metricData.TargetDefectStatPo) {
	return func(key string, value *metricData.TargetDefectStatPo) {
		if targetDefectStat, ok := (*targetDefectStats)[key]; ok {
			targetDefectStat.DefectNum += 1
		} else {
			(*targetDefectStats)[key] = value
		}
	}
}

func saveTargetDefectStatRepository(targetDefectStats map[string]*metricData.TargetDefectStatPo) {
	var defectStats []metricData.TargetDefectStatPo
	for _, value := range targetDefectStats {
		defectStats = append(defectStats, *value)
	}
	err := metricRepo.NewTargetDefectStatRepository().SaveInBatch(&defectStats)
	if err != nil {
		logger.Errorf("save target defect statistic: %s", err.Error())
		return
	}
}

func updateTargetDefectStatRepository(targetDefectStats map[string]*metricData.TargetDefectStatPo) {
	var defectStats []metricData.TargetDefectStatPo
	for _, value := range targetDefectStats {
		defectStats = append(defectStats, *value)
	}
	err := metricRepo.NewTargetDefectStatRepository().UpdateInBatch(&defectStats)
	if err != nil {
		logger.Errorf("update target defect statistic: %s", err.Error())
		return
	}
}

type response struct {
	Result string `json:"result"`
	Code   int    `json:"code"`
	Status string `json:"status"`
	Detail Detail `json:"detail"`
}

type Detail struct {
	Bo    interface{} `json:"bo"`
	Code  interface{} `json:"code"`
	Other Other       `json:"other"`
}

type Other struct {
	Message string `json:"message"`
}

func SyncTargetDefectStat() {
	err := syncDefectToDatabase()
	if err != nil {
		logger.Errorf("sync all defect to database failed, Message: %s", err.Error())
		return
	}
	targets, err := target.GetTargetService().FindAllTarget()
	if err != nil {
		logger.Errorf("find all target failed, Message: %s", err.Error())
		return
	}
	err = metricRepo.NewTargetDefectStatRepository().RemoveAll()
	if err != nil {
		logger.Errorf("delete all target_defect_stat failed, Message: %s", err.Error())
		return
	}
	for i, _ := range targets {
		tar := &entity.TargetEntity{Id: targets[i].Id}
		answerFile := tar.GetAnswerFileDir()
		fileName := filepath.Join(answerFile, targets[i].Answer)
		answerShootingData, err := shootingresult.NewShootingResultCalculator().LoadShootingData(fileName)
		if err != nil {
			logger.Errorf("err %v", err)
			continue
		}
		targetDefectStats := make(map[string]*metricData.TargetDefectStatPo, 10)
		saveTargetDefect := saveTargetDefectStat(&targetDefectStats)

		for _, data := range answerShootingData {
			saveTargetDefect(data.DefectCode, &metricData.TargetDefectStatPo{
				ID:        uuid.NewString(),
				TargetID:  tar.Id,
				DefectId:  data.DefectCode,
				DefectNum: 1,
			})
		}
		saveTargetDefectStatRepository(targetDefectStats)
	}
}

// todoGetCurrentTemplateFileName
func syncDefectToDatabase() error {
	fileName, _ := controller.GetCurrentTemplateFileName(util.DefaultWorkspace, "")
	err := defect.UpdateDefectCoder(filepath.Join(util.TemplateDir, util.DefaultWorkspace, fileName))
	if err != nil {
		return errors.WithMessagef(err, "sync defects to database failed")
	}
	return nil
}
