package document

import (
	"code-shooting/domain/entity"
	"code-shooting/infra/util"
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"
)

const (
	ManagerDir string = "system-manage"
	UserDir    string = "user-guide"

	UserType    string = "user"
	ManagerType string = "manager"
	AllType     string = "all"
)

var (
	rootPath    string = util.ConfDir + "/doc"
	userPath    string = util.ConfDir + "/doc/user-guide"
	managerPath string = util.ConfDir + "/doc/system-manage"
)

type DocumentService struct {}

func NewDocumentService() *DocumentService{
	return &DocumentService{}
}

func (this *DocumentService)GetDocumentDir(docType string) entity.DocumentsInfo {
	switch docType {
	case AllType:
		return findDocumentDir(rootPath, "", docType)
	case ManagerType:
		return findDocumentDir(managerPath, "", docType)
	case UserType:
		return findDocumentDir(userPath, "", docType)
	}
	return entity.DocumentsInfo{}
}

func findDocumentDir(filePath string, parentID string, docType string) entity.DocumentsInfo {
	result := entity.DocumentsInfo{
		ID:       filePath,
		Name:     filepath.Base(filePath),
		ParentId: parentID,
		Type:     docType,
	}

	fileInfos, err := ioutil.ReadDir(filePath)
	if err != nil {
		return result
	}
	for _, fileInfo := range fileInfos {
		if !fileInfo.IsDir() {
			result.Ducoments = append(result.Ducoments, fileInfo.Name())
			continue
		}
		tempType := docType
		if fileInfo.Name() == ManagerDir {
			tempType = "manager"
		}
		if fileInfo.Name() == UserDir {
			tempType = "user"
		}
		result.Children = append(result.Children, findDocumentDir(fmt.Sprintf("%s/%s", filePath, fileInfo.Name()), result.Name, tempType))
	}
	return result
}

func (this *DocumentService) GetDocumentDetail(catalogueId, filePath string) string {
	isValid := strings.HasSuffix(filePath, ".md")
	if !isValid {
		return ""
	}
	content, err := ioutil.ReadFile(path.Clean(catalogueId + "/" + filePath))
	if err != nil {
		return ""
	}
	return string(content)
}
