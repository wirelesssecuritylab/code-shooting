package shootingresult

import (
	"code-shooting/infra/shooting-result/defect"
	"code-shooting/infra/shooting-result/util"
	"strconv"

	"code-shooting/infra/logger"

	"github.com/pkg/errors"
)

const (
	record = "记录"
)

const (
	fName int = iota
	sLine
	eLine
	dClass
	dSubClass
	dDescribe
)

var titles = []string{
	"文件名",
	"起始行号",
	"结束行号",
	"缺陷大类",
	"缺陷小类",
	"缺陷细项",
}

type TargetAnswer struct {
	FileName       string
	StartLineNum   int
	EndLineNum     int
	StartColNum    int
	EndColNum      int
	DefectClass    string
	DefectSubClass string
	DefectDescribe string
	DefectCode     string
	Remark         string
}

func (s *ResultCalculator) LoadShootingData(file string) ([]TargetAnswer, error) {
	sheetRows, err := s.getSheetRowsByName(record, file)
	if err != nil {
		return nil, errors.Wrapf(err, "get target record sheet rows")
	}

	indexes := util.GetIndexesByNames(titles, sheetRows)
	if !util.IsIndexesValid(indexes, titles) {
		return nil, errors.Errorf("target record sheet title is wrong")
	}
	coder, err := defect.NewDefectEncoder(file)
	if err != nil {
		return nil, errors.Wrapf(err, "new defect coder")
	}

	shootResult := NewShootingResultCalculator()
	language := shootResult.getLanguage(file)
	if language == "C" || language == "C++" {
		language = "C&C++"
	}
	var targetAnswers []TargetAnswer
	for _, row := range sheetRows[1:] {
		lenOfRow := len(row)
		if !util.IsRowValid(indexes, lenOfRow) {
			logger.Warnf("invalid target ring, it will be ignored")
			continue
		}
		snum, err := strconv.Atoi(row[indexes[sLine]])
		if err != nil {
			logger.Warnf("invalid target ring, startline number is not an int number, it will be ignored")
			continue
		}
		enum, err := strconv.Atoi(row[indexes[eLine]])
		if err != nil {
			logger.Warnf("invalid target ring, endline number is not an int number, it will be ignored")
			continue
		}

		defectCode := coder.EncodeDefect(row[indexes[dClass]], row[indexes[dSubClass]], row[indexes[dDescribe]])
		if defectCode == "" {
			defectCode = coder.EncodeDefect(row[indexes[dClass]], row[indexes[dSubClass]], row[indexes[dDescribe]])
		}
		targetAnswers = append(targetAnswers,
			TargetAnswer{
				FileName:       row[indexes[fName]],
				StartLineNum:   snum,
				EndLineNum:     enum,
				DefectClass:    row[indexes[dClass]],
				DefectSubClass: row[indexes[dSubClass]],
				DefectDescribe: row[indexes[dDescribe]],
				DefectCode:     defectCode,
			},
		)
	}

	return targetAnswers, nil
}
