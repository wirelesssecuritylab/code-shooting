package shootingresult

import (
	"strconv"

	"strings"

	"code-shooting/infra/logger"

	"github.com/pkg/errors"
	"github.com/xuri/excelize/v2"

	"code-shooting/domain/entity"
	"code-shooting/domain/service/score"
	"code-shooting/infra/shooting-result/defect"
)

type check func(src *score.TargetAnswer, dst *TargetAnswer) bool
type ResultCalculator struct{}

func NewShootingResultCalculator() *ResultCalculator {
	return &ResultCalculator{}
}

func (s *ResultCalculator) getAnswerRingAndScoreInfo(answerFile string) ([]TargetAnswer, map[string]int, error) {
	coder, err := defect.NewDefectEncoder(answerFile) // 获取靶标文件中所有映射：{(缺陷大类，小类，细项) ： 缺陷编码}
	if err != nil {
		return nil, nil, errors.Wrapf(err, "new defect coder")
	}

	answers, err := s.LoadShootingData(answerFile) // 获取靶标文件中靶环信息列表
	if err != nil {
		return nil, nil, errors.Wrapf(err, "load answer shooting data")
	}
	for i := range answers {
		t := &answers[i]
		answers[i].DefectCode = coder.EncodeDefect(t.DefectClass, t.DefectSubClass, t.DefectDescribe)
	}

	scoreCfg, err := s.LoadScoreConfig(answerFile) // 获取当前靶子对应语言和通用缺陷分类中缺陷大类所对应分值{"缺陷大类"：分值}
	if err != nil {
		return nil, nil, errors.Wrapf(err, "load score config")
	}

	return answers, scoreCfg, nil
}

func (s *ResultCalculator) GetAnswerRingNumAndScore(answerFile string) (num, score uint32, err error) {
	answers, scoreCfg, err := s.getAnswerRingAndScoreInfo(answerFile)
	if err != nil {
		return 0, 0, errors.Wrapf(err, "parse answer rings and score cfg")
	}
	num = uint32(len(answers))
	score = uint32(s.getTotalPointOfAnswer(answers, scoreCfg))
	err = nil
	return
}

func (s *ResultCalculator) CalculateShootingResult(answerFile string, tp *score.TargetAnswerPaper) (*entity.TargetResult, error) {
	targetAnswers, scoreCfg, err := s.getAnswerRingAndScoreInfo(answerFile)
	if err != nil {
		return nil, errors.Wrapf(err, "parse answer rings and score cfg")
	}

	res := &entity.TargetResult{
		TargetDetails: make([]entity.TargetDetail, 0),
	}

	num := 0
	point := 0
	total_num := len(targetAnswers)
	total_point := s.getTotalPointOfAnswer(targetAnswers, scoreCfg)
	for _, sd := range tp.Answers { //用户答卷列表
		isMatch := false
		index := -1
		for i, asd := range targetAnswers { //靶环列表
			if s.isTargetMatch(&sd, &asd) {
				isMatch = true
				index = i
				break
			}
		}
		target := entity.TargetDetail{
			FileName:     sd.FileName,
			StartLineNum: sd.StartLineNum,
			EndLineNum:   sd.EndLineNum,
			StartColNum:  sd.StartColNum,
			EndColNum:    sd.EndColNum,
			DefectCode:   sd.DefectCode,
			Remark:       sd.Remark,
			TargetScore:  0,
		}
		if isMatch {
			num += 1
			asd := targetAnswers[index]
			point += scoreCfg[asd.DefectClass]
			targetAnswers = append(targetAnswers[:index], targetAnswers[index+1:]...)
			target.TargetScore = scoreCfg[asd.DefectClass]
		}
		res.TargetDetails = append(res.TargetDetails, target)
	}
	res.HitNum = num
	res.TotalNum = total_num
	res.HitScore = point
	res.TotalScore = total_point
	return res, nil
}

func (s *ResultCalculator) isTargetMatch(src *score.TargetAnswer, dst *TargetAnswer) bool {
	checkers := []check{
		s.checkFileName,
		s.checkDefectCode,
		s.checkLineNum, //顺序必须放在最后，这样才能实现靶标0，0匹配所有行
	}
	for _, check := range checkers {
		if !check(src, dst) {
			return false
		}
	}
	return true
}

func (s *ResultCalculator) getTotalPointOfAnswer(targetAnswers []TargetAnswer, scoreCfg map[string]int) int {
	point := 0
	for _, sd := range targetAnswers {
		point += scoreCfg[sd.DefectClass]
	}
	return point
}

func (s *ResultCalculator) LoadScoreConfig(answerFile string) (map[string]int, error) {
	cfg := make(map[string]int)
	language := s.getLanguage(answerFile)
	if language != "" {
		lngRows, err := s.getSheetRowsByName(language, answerFile)
		if err != nil {
			return nil, errors.Wrapf(err, "get %s defect class sheet rows", language)
		}
		s.parseDefectClassRows(lngRows, cfg)
	}

	generalRows, err := s.getSheetRowsByName("通用", answerFile)
	if err != nil {
		logger.Warnf("get generic defect class sheet rows")
		return cfg, nil
	} else {
		s.parseDefectClassRows(generalRows, cfg)
	}

	return cfg, nil
}

func (s *ResultCalculator) parseDefectClassRows(sheetRows [][]string, scoreCfg map[string]int) {
	if len(sheetRows) <= 0 {
		return
	}
	scoreIndex := -1
	for i := range sheetRows[0] {
		if sheetRows[0][i] == "分值" {
			scoreIndex = i
			break
		}
	}
	if scoreIndex == -1 {
		return
	}
	for _, row := range sheetRows[1:] {
		if len(row) <= scoreIndex {
			continue
		}
		score, err := strconv.Atoi(row[scoreIndex])
		if err != nil {
			continue
		}
		if _, ok := scoreCfg[row[0]]; !ok || score > scoreCfg[row[0]] {
			scoreCfg[row[0]] = score
		}
	}
}

func (s *ResultCalculator) getLanguage(answerFile string) string {
	sheetRows, err := s.getSheetRowsByName(record, answerFile)
	if err != nil {
		return ""
	}
	if len(sheetRows) <= 1 {
		return ""
	}
	var index = -1
	for i := range sheetRows[0] {
		if sheetRows[0][i] == "语言类型" {
			index = i
			break
		}
	}

	if len(sheetRows[1]) <= index || index == -1 {
		return ""
	}

	return sheetRows[1][index]
}

func (s *ResultCalculator) getSheetRowsByName(substr, file string) ([][]string, error) {
	f, err := excelize.OpenFile(file)
	if err != nil {
		return nil, errors.Wrapf(err, "open excel file %s", file)
	}
	defer f.Close()
	sheetName := ""
	for _, sheet := range f.WorkBook.Sheets.Sheet {
		if strings.Contains(sheet.Name, substr) {
			sheetName = sheet.Name
			break
		}
	}
	if sheetName == "" {
		return nil, errors.Errorf("can't find sheet")
	}
	sheetRows, err := f.GetRows(sheetName)
	if err != nil || len(sheetRows) == 0 {
		return nil, errors.Errorf("sheet is empty with error: %v", err)
	}
	return sheetRows, nil
}

func (s *ResultCalculator) checkFileName(src *score.TargetAnswer, dst *TargetAnswer) bool {
	return strings.TrimSpace(src.FileName) == strings.TrimSpace(dst.FileName)
}

func (s *ResultCalculator) checkLineNum(src *score.TargetAnswer, dst *TargetAnswer) bool {
	if dst.EndLineNum == 0 && dst.StartLineNum == 0 {
		return true
	}

	if src.StartLineNum < dst.StartLineNum {
		return false
	}
	if src.EndLineNum > dst.EndLineNum {
		return false
	}

	return true
}

func (s *ResultCalculator) checkDefectCode(src *score.TargetAnswer, dst *TargetAnswer) bool {
	return src.DefectCode == dst.DefectCode
}
