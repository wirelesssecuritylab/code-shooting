package defect

import (
	"code-shooting/infra/logger"
	"code-shooting/infra/po"
	"code-shooting/infra/repository"
	"code-shooting/infra/shooting-result/util"
	"strings"

	"github.com/pkg/errors"
	"github.com/xuri/excelize/v2"
)

func UpdateDefectCoder(filePath string) error {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return errors.Wrapf(err, "open excel file %s", filePath)
	}
	defer f.Close()

	defectRep := repository.NewDefectRepository()
	err = defectRep.RemoveAll()
	if err != nil {
		return errors.Wrapf(err, "delete all defect failed")
	}
	var defectPos []po.DefectPo
	for _, sheet := range f.WorkBook.Sheets.Sheet {
		if !strings.Contains(sheet.Name, "缺陷分类") {
			continue
		}
		sheetRows, err := f.GetRows(sheet.Name)
		if err != nil {
			return errors.Wrapf(err, "get sheet %s rows", sheet.Name)
		}
		defects := getSheetDefects(sheetRows)
		language := getLanguage(sheet.Name)
		class := ""
		subClass := ""
		for _, defect := range defects {
			if defect[dClassIndex] != "" {
				class = defect[dClassIndex]
			}
			if defect[dSubClassIndex] != "" {
				subClass = defect[dSubClassIndex]
			}
			flag := false
			for i, _ := range defectPos {
				if defectPos[i].DefectId == defect[dCodeIndex] {
					flag = true
					logger.Warnf("find duplicate defect code: %s", defect[dCodeIndex])
					break
				}
			}
			if flag {
				continue
			}
			defectPos = append(defectPos, po.DefectPo{
				DefectId:       defect[dCodeIndex],
				DefectClass:    class,
				DefectSubclass: subClass,
				DefectDescribe: defect[dDescribeIndex],
				Language:       language,
			})
		}
	}
	err = defectRep.SaveInBatch(&defectPos)
	if err != nil {
		return errors.Wrapf(err, "save all defect failed")
	}
	return nil
}

func getSheetDefects(rows [][]string) [][]string {
	indexes := util.GetIndexesByNames(titles, rows)
	if !util.IsIndexesValid(indexes, titles) {
		return [][]string{}
	}

	var defects [][]string
	for _, row := range rows[1:] {
		if !util.IsRowValid(indexes, len(row)) {
			continue
		}

		defect := []string{
			row[indexes[dClassIndex]],
			row[indexes[dSubClassIndex]],
			strings.Replace(row[indexes[dDescribeIndex]], ",", "，", -1), // 解决无法匹配半角逗号
			row[indexes[dCodeIndex]],
		}
		defects = append(defects, defect)
	}
	return defects
}

func getLanguage(sheetName string) string {
	sheetName = strings.Replace(sheetName, "（", "(", -1)
	sheetName = strings.Replace(sheetName, "）", ")", -1)
	containLanguage := strings.Split(sheetName, "(")
	if len(containLanguage) < 2 {
		return ""
	}
	endIndex := strings.Index(containLanguage[1], ")")
	language := containLanguage[1][:endIndex]
	return language
}
