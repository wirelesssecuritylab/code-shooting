package defect

import (
	"code-shooting/domain/service/score"
	"code-shooting/infra/shooting-result/util"
	"encoding/base64"
	"fmt"
	"io"
	"strings"

	"github.com/pkg/errors"
	"github.com/xuri/excelize/v2"
)

var titles = []string{
	"缺陷大类",
	"缺陷小类",
	"缺陷细项",
	"缺陷编码",
	"操作说明",
}

const (
	dClassIndex int = iota
	dSubClassIndex
	dDescribeIndex
	dCodeIndex
	dOperationIndex
)

type defectCoder struct {
	defectsToId map[string]string
	idToDefects map[string]string
}

func NewDefectEncoder(filePath string) (score.DefectCoder, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, errors.Wrapf(err, "open excel file %s", filePath)
	}
	defer f.Close()

	coder := &defectCoder{
		defectsToId: make(map[string]string),
		idToDefects: make(map[string]string),
	}
	for _, sheet := range f.WorkBook.Sheets.Sheet {
		if !strings.Contains(sheet.Name, "缺陷分类") {
			continue
		}
		sheetRows, err := f.GetRows(sheet.Name)
		if err != nil {
			return nil, errors.Wrapf(err, "get sheet %s rows", sheet.Name)
		}
		defects := coder.getSheetDefects(sheetRows)

		class := ""
		subClass := ""
		for _, defect := range defects {
			if defect[dClassIndex] != "" {
				class = defect[dClassIndex]
			}
			if defect[dSubClassIndex] != "" {
				subClass = defect[dSubClassIndex]
			}
			if defect[dOperationIndex] == moveOperation || defect[dOperationIndex] == discardOperation {
				continue
			}
			defectKey := coder.getDefect(class, subClass, defect[dDescribeIndex])
			coder.defectsToId[defectKey] = defect[dCodeIndex]
			coder.idToDefects[defect[dCodeIndex]] = defectKey
		}

	}

	return coder, nil
}

func (s *defectCoder) EncodeDefect(dClass, dSubClass, dDescribe string) string {
	key := s.getDefect(dClass, dSubClass, dDescribe)
	id, ok := s.defectsToId[key]
	if ok {
		return id
	}
	return ""
}

func (s *defectCoder) DecodeDefect(id string) (string, string, string) {
	defect, ok := s.idToDefects[id]
	if !ok {
		return "", "", ""
	}

	return s.parseDefect(defect)
}

func (s *defectCoder) getDefect(dClass, dSubClass, dDescribe string) string {
	return base64.StdEncoding.EncodeToString([]byte(dClass)) + "-" +
		base64.StdEncoding.EncodeToString([]byte(dSubClass)) + "-" + base64.StdEncoding.EncodeToString([]byte(dDescribe))
}

func (s *defectCoder) parseDefect(defect string) (string, string, string) {
	defectInfo := strings.Split(defect, "-")
	if len(defectInfo) != 3 {
		return "", "", ""
	}

	dClass, err := base64.StdEncoding.DecodeString(defectInfo[0])
	if err != nil {
		return "", "", ""
	}
	dSubClass, err := base64.StdEncoding.DecodeString(defectInfo[1])
	if err != nil {
		return "", "", ""
	}
	dDescribe, err := base64.StdEncoding.DecodeString(defectInfo[2])
	if err != nil {
		return "", "", ""
	}
	return string(dClass), string(dSubClass), string(dDescribe)
}

func (s *defectCoder) getSheetDefects(rows [][]string) [][]string {
	indexes := util.GetIndexesByNames(titles, rows)
	newTitles := titles
	if indexes[dOperationIndex] == -1 {
		indexes = indexes[:len(indexes)-1]
		newTitles = newTitles[:len(newTitles)-1]
	}
	if !util.IsIndexesValid(indexes, titles) {
		return [][]string{}
	}

	var defects [][]string
	for _, row := range rows[1:] {
		if !util.IsRowValid(indexes[:len(indexes)-1], len(row)) {
			continue
		}

		operation := ""
		if len(newTitles) == len(titles) && indexes[dOperationIndex] < len(row) {
			operation = row[indexes[dOperationIndex]]
		}

		defect := []string{
			row[indexes[dClassIndex]],
			row[indexes[dSubClassIndex]],
			strings.Replace(row[indexes[dDescribeIndex]], ",", "，", -1), // 解决无法匹配半角逗号
			row[indexes[dCodeIndex]],
			operation,
		}
		defects = append(defects, defect)
	}
	return defects
}

const (
	moveOperation    = "move"
	discardOperation = "discard"
)

var (
	defectInfoIndex    = 0
	defectCodeIndex    = 1
	operationInfoIndex = 2
)

func CheckDefectCode(file io.Reader) error {
	f, err := openExcelFromReader(file)
	if err != nil {
		return errors.WithMessage(err, "openExcelFromMultipartFile")
	}
	defer f.Close()

	m := make(map[string]string)

	for _, sheet := range f.WorkBook.Sheets.Sheet {
		if !strings.Contains(sheet.Name, "缺陷分类") {
			continue
		}
		sheetRows, err := f.GetRows(sheet.Name)
		if err != nil {
			return errors.WithMessagef(err, "get sheet %s rows", sheet.Name)
		}
		titles := []string{"缺陷细项", "缺陷编码", "操作说明"}
		indexes := util.GetIndexesByNames(titles, sheetRows)
		if indexes[operationInfoIndex] == -1 {
			indexes = indexes[:len(indexes)-1]
		}
		if !isIndexesValid(indexes) {
			return errors.New("规范格式错误")
		}
		for _, sheetRow := range sheetRows[1:] {
			if len(indexes) == len(titles) && util.IsRowValid(indexes, len(sheetRow)) && sheetRow[indexes[operationInfoIndex]] == moveOperation {
				continue
			}
			if _, ok := m[sheetRow[indexes[defectCodeIndex]]]; ok {
				return errors.New(fmt.Sprintf("缺陷编码%s-%s重复", sheet.Name, sheetRow[indexes[defectCodeIndex]]))
			}
			m[sheetRow[indexes[defectCodeIndex]]] = sheetRow[indexes[defectInfoIndex]]
		}
	}

	return nil
}
func isIndexesValid(index []int) bool {
	for _, i := range index {
		if i == -1 {
			return false
		}
	}
	return true
}
func openExcelFromReader(reader io.Reader) (*excelize.File, error) {
	f, err := excelize.OpenReader(reader)
	if err != nil {
		return f, err
	}
	f.Path = ""
	return f, nil
}
