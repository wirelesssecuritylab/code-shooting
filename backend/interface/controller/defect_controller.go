package controller

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"code-shooting/infra/logger"
	"code-shooting/infra/restserver"

	"github.com/xuri/excelize/v2"

	"code-shooting/app/service/template"
	"code-shooting/domain/service/target"
	"code-shooting/infra/common"
	"code-shooting/infra/util"
	"code-shooting/interface/assembler"
	"code-shooting/interface/dto"
)

const (
	defectCommonSheetName = "缺陷分类（通用）"
	sheetNamePrefix       = "缺陷分类"
)

type DefectController struct {
	Result string      `json:"result"`
	Detail interface{} `json:"detail"`
}

var langSheetMap = map[string]string{
	"go":         "缺陷分类（Go）",
	"c&c++":      "缺陷分类（C&C++）",
	"c":          "缺陷分类（C&C++）",
	"c++":        "缺陷分类（C&C++）",
	"python":     "缺陷分类（Python）",
	"java":       "缺陷分类（Java）",
	"scala":      "缺陷分类（Scala）",
	"typescript": "缺陷分类（TypeScript）",
	"javascript": "缺陷分类（JavaScript）",
}

func NewDefectController() *DefectController {
	return &DefectController{}
}

func (s *DefectController) Post(ctx restserver.Context) error {
	body, _ := ioutil.ReadAll(ctx.Request().Body)

	var request = &dto.DefectAction{}
	err := assembler.ParseReq(body, request)
	if err != nil {
		return ctx.String(http.StatusBadRequest, fmt.Sprintf("bad request: %s", err.Error()))
	}

	if defects, err := DoGetDefects(ctx.Param("id"), request.Lang, request.NeedCode); err != nil {
		logger.Error("get defect failed of : ", err)
		return ctx.JSON(http.StatusNotFound, &common.Response{Result: "failed"})
	} else {
		defectResult := DefectController{"success", defects}
		return ctx.JSON(http.StatusOK, defectResult)
	}
}

/*
 * 通过靶子ID获取任务空间，然后根据工作空间和语言查询EXCEL中具体的缺陷信息
 */
func DoGetDefects(targetID, language string, needCode bool) (dto.Defects, error) {
	workspace := ""
	t, err := target.GetTargetService().FindByID(targetID)
	if err == nil {
		workspace = t.Workspace
	}
	if workspace == "" {
		workspace = util.DefaultWorkspace
	}
	//传入靶子模板版本
	sheet, err := getSheetByWorkspace(workspace, language, t.Template)
	if err != nil {
		return nil, err
	}

	return ReadSheetDefects(workspace, sheet, t.Template, needCode)
}

func ReadSheetDefects(workspace, sheet string, templateVersion string, needCode bool) (dto.Defects, error) {
	fileName, _ := GetCurrentTemplateFileName(workspace, templateVersion)
	if fileName == "" {
		return nil, fmt.Errorf("get current template file failed from workspace %s", workspace)
	}

	f, err := excelize.OpenFile(filepath.Join(util.TemplateDir, workspace, fileName))
	if err != nil {
		return nil, fmt.Errorf("open excel file %s failed: %s ", filepath.Join(util.TemplateDir, workspace, fileName), err.Error())
	}
	defer f.Close()
	m := make(dto.Defects)

	getDefectDetail(f, defectCommonSheetName, m, needCode)
	getDefectDetail(f, sheet, m, needCode)

	return m, nil
}

func getDefectDetail(f *excelize.File, sheet string, detail dto.Defects, needCode bool) error {
	defectIsObsolete := func(currentLine int, defectBig, defectSmall, defectDetail string) bool {
		if columnName, _ := f.GetCellValue(sheet, "E1"); strings.EqualFold(strings.Trim(columnName, " "), "操作说明") {
			operate, _ := f.GetCellValue(sheet, fmt.Sprintf("E%d", currentLine))
			if strings.EqualFold(strings.Trim(operate, " "), "move") || strings.EqualFold(strings.Trim(operate, " "), "discard") {
				logger.Info(fmt.Sprintf("C%d", currentLine), " is "+operate+" defect of ", defectBig, " ", defectSmall, " ", defectDetail)
				return true
			}

			return false
		}

		cellStyleID, _ := f.GetCellStyle(sheet, fmt.Sprintf("C%d", currentLine))
		obsoleteIDs := []int{107, 99}
		for _, obsoleteID := range obsoleteIDs {
			if cellStyleID == obsoleteID {
				logger.Info(fmt.Sprintf("C%d", currentLine), " is obsolete defect of ", defectBig, " ", defectSmall, " ", defectDetail, " style ID is ", cellStyleID)
				return true
			}
		}
		return false
	}

	line := 2
	for {
		defectBig, _ := f.GetCellValue(sheet, fmt.Sprintf("A%d", line))
		defectSmall, _ := f.GetCellValue(sheet, fmt.Sprintf("B%d", line))
		defectDetail, _ := f.GetCellValue(sheet, fmt.Sprintf("C%d", line))

		defectCode := ""
		defectDetail = strings.Replace(defectDetail, ",", "，", -1) // 解决无法匹配半角逗号

		if needCode {
			defectCode, _ = f.GetCellValue(sheet, fmt.Sprintf("D%d", line))
		}
		if defectBig == "" {
			break
		}

		line += 1

		if defectIsObsolete(line-1, defectBig, defectSmall, defectDetail) {
			continue
		}

		if _, ok := detail[defectBig]; ok {
			if _, ok := detail[defectBig][defectSmall]; ok {
				detail[defectBig][defectSmall] = append(detail[defectBig][defectSmall], dto.DefectDetail{
					Description: defectDetail,
					Code:        defectCode})
			} else {
				defectDetails := []dto.DefectDetail{{
					Description: defectDetail,
					Code:        defectCode}}
				detail[defectBig][defectSmall] = defectDetails
			}
		} else {
			defectDetails := []dto.DefectDetail{{
				Description: defectDetail,
				Code:        defectCode}}
			detailSmalls := make(map[string][]dto.DefectDetail)
			detailSmalls[defectSmall] = defectDetails
			detail[defectBig] = detailSmalls
		}
	}

	return nil
}

func GetCurrentTemplateFileName(workspace string, templateVersion string) (string, error) {
	if workspace == "" {
		workspace = util.DefaultWorkspace
	}
	templates, err := template.GetTemplateAppService().QueryTemplate()
	if err != nil {
		return "", err
	}
	if templateVersion != "" {
		for index := range templates {
			if templates[index].Active && templates[index].Workspace == workspace && templates[index].Version == templateVersion {
				return template.BuildTemplateFileName(templates[index].Version), nil
			}
		}
	} else {
		var versionsActive []string
		for index := range templates {
			if templates[index].Active && templates[index].Workspace == workspace {
				versionsActive = append(versionsActive, (templates[index].Version)[1:len((templates[index].Version))])
			}
		}
		if len(versionsActive) == 0 {
			return "", fmt.Errorf("工作空间 %s 中不存在启用中规范", workspace)
		}
		maxVersion := versionsActive[0]
		maxVersion_a, _ := strconv.Atoi(strings.Split(maxVersion, ".")[0])
		maxVersion_b, _ := strconv.Atoi(strings.Split(maxVersion, ".")[1])
		for i := 0; i < len(versionsActive); i++ {
			version_a, _ := strconv.Atoi(strings.Split(versionsActive[i], ".")[0])
			version_b, _ := strconv.Atoi(strings.Split(versionsActive[i], ".")[1])
			if maxVersion_a < version_a {
				maxVersion_a = version_a
				maxVersion_b = version_b
			} else if maxVersion_a == version_a {
				if maxVersion_b < version_b {
					maxVersion_b = version_b
				}
			}
		}
		return template.BuildTemplateFileName("v" + strconv.Itoa(maxVersion_a) + "." + strconv.Itoa(maxVersion_b)), nil
	}
	return "", fmt.Errorf("工作空间 %s 中不存在启用中规范", workspace)
}

func GetWorkSpaceSheets(workspace string, templateVersion string) ([]string, error) {
	fileName, _ := GetCurrentTemplateFileName(workspace, templateVersion)
	if fileName == "" {
		return nil, fmt.Errorf("get current template file failed from workspace %s", workspace)
	}

	f, err := excelize.OpenFile(filepath.Join(util.TemplateDir, workspace, fileName))
	if err != nil {
		return nil, fmt.Errorf("open excel file %s failed: %s ", filepath.Join(util.TemplateDir, workspace, fileName), err.Error())
	}

	defer f.Close()

	sheets := f.GetSheetList()
	return sheets, nil
}

/*
*

	优化，根据启用的多规范获取语言
*/
func (s *DefectController) LoadSheetLangUpdate(ctx restserver.Context) error {
	workspace := "public"
	if "" != ctx.Param("workspace") {
		workspace = ctx.Param("workspace")
	}
	templates, err := template.GetTemplateAppService().QueryTemplate()
	if err != nil {
		return err
	}
	langsReturn := make(map[string]int)
	for i := 0; i < len(templates); i++ {
		if (workspace == templates[i].Workspace) && (templates[i].Active == true) {
			fileName := template.BuildTemplateFileName(templates[i].Version)
			if fileName == "" {
				return fmt.Errorf("get current template file failed from workspace %s", workspace)
			}
			f, err := excelize.OpenFile(filepath.Join(util.TemplateDir, workspace, fileName))
			if err != nil {
				return fmt.Errorf("open excel file %s failed: %s ", filepath.Join(util.TemplateDir, workspace, fileName), err.Error())
			}
			defer f.Close()
			sheets := f.GetSheetList()
			langs := formatSheetNames(sheets)
			for k, v := range langs {
				langsReturn[k] = v
			}
		}
	}
	return ctx.JSON(http.StatusOK, langsReturn)
}
func (s *DefectController) GetLangsByWorkspaceAndTemplate(ctx restserver.Context) error {
	workspace := ctx.Param("workspace")
	templateVersion := ctx.Param("template")
	sheets, err := GetWorkSpaceSheets(workspace, templateVersion)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, fmt.Sprintf("load failed : %s", err.Error()))
	}
	langs := formatSheetNames(sheets)
	langsReturn := make(map[string]int)
	for k, v := range langs {
		langsReturn[k] = v
	}
	return ctx.JSON(http.StatusOK, langsReturn)
}

/*
*
根据工作空间以及语言获取打靶规范
*/
func (s *DefectController) GetTemplateByWorksapce(ctx restserver.Context) error {
	workspace := ctx.Param("workspace")
	templatesVersion, err := getTemplateVersion(workspace)
	if err != nil {
		return ctx.String(http.StatusNotFound, fmt.Sprintf("load failed : %s", err.Error()))
	}

	return ctx.JSON(http.StatusOK, templatesVersion)
}

func getTemplateVersion(workspace string) ([]string, error) {
	var versionSlice []string
	if workspace == "" {
		workspace = util.DefaultWorkspace
	}
	templates, err := template.GetTemplateAppService().QueryTemplate()
	if err != nil {
		return versionSlice, err
	}
	for i := 0; i < len(templates); i++ {
		if (workspace == templates[i].Workspace) && (templates[i].Active == true) {
			versionSlice = append(versionSlice, templates[i].Version)
		}
	}
	return versionSlice, nil
}

func (s *DefectController) GetLangDefect(ctx restserver.Context) error {
	workspace := ctx.Param("workspace")
	body, _ := ioutil.ReadAll(ctx.Request().Body)

	var request = &dto.DefectAction{}
	err := assembler.ParseReq(body, request)
	if err != nil {
		return ctx.String(http.StatusBadRequest, fmt.Sprintf("bad request: %s", err.Error()))
	}
	sheet, err := getSheetByWorkspace(workspace, request.Lang, request.TemplateVersion)
	if err != nil {
		return ctx.String(http.StatusNotFound, fmt.Sprintf("load failed : %s", err.Error()))
	}
	if defects, err := ReadSheetDefects(workspace, sheet, request.TemplateVersion, request.NeedCode); err != nil {
		logger.Error("get defect failed of : ", err)
		return ctx.JSON(http.StatusNotFound, &common.Response{Result: "failed"})
	} else {
		defectResult := DefectController{"success", defects}
		return ctx.JSON(http.StatusOK, defectResult)
	}
}

func trimRule(s rune) bool {
	if s == '（' || s == '）' {
		return true
	}
	return false
}

/*
通过工作空间和语言查询Excel中具体的sheet名称
*/
func getSheetByWorkspace(workspace string, language string, templateVersion string) (string, error) {
	sheet := ""
	sheets, err := GetWorkSpaceSheets(workspace, templateVersion)
	if err != nil {
		return sheet, fmt.Errorf("load failed : %s", err.Error())
	}
	langs := formatSheetNames(sheets)
	for key, value := range langs {
		// 通过格式化的语言类型去索引sheet名，与前端保持一致，避免大小，空格等因素导致的异常
		if key == language {
			sheet = sheets[value]
			return sheet, nil
		}
	}
	return sheet, fmt.Errorf("not found language %s from workspace %s", language, workspace)
}

/*
提取sheet名称中的括号的内容
*/
func formatSheetNames(names []string) map[string]int {

	var lang []string
	langInfo := make(map[string]int)

	for index, value := range names {
		if strings.Contains(value, sheetNamePrefix) && value != defectCommonSheetName {
			val := strings.TrimPrefix(value, sheetNamePrefix)
			val = strings.TrimFunc(val, trimRule)
			val = strings.TrimSpace(val)
			lang = append(lang, val)
			langs := strings.Split(val, "&")
			for _, l := range langs {
				langInfo[l] = index
			}
		}
	}
	return langInfo
}
