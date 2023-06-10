package target

import (
	"archive/zip"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"code-shooting/infra/logger"

	"github.com/pkg/errors"

	"code-shooting/domain/entity"
	repo "code-shooting/domain/repository"
	ecsvc "code-shooting/domain/service/ec"
	"code-shooting/infra/errcode"
	shootingresult "code-shooting/infra/shooting-result"
	sr "code-shooting/infra/shooting-result"
	"code-shooting/infra/util/tools"
	"code-shooting/interface/dto"
	"sort"
	"time"
)

type TargetService struct {
	repo.TargetInterface
}

var targetService *TargetService

func GetTargetService() *TargetService {
	if targetService == nil {
		targetService = &TargetService{repo.GetTargetRepo()}
	}
	return targetService
}

func (s *TargetService) CreateTarget(t *dto.TargetModel) error {
	te := NewTarget(t)
	if err := s.InsertTarget(te); err != nil {
		return err
	}
	if err := MoveTmpFiles2Target(te); err != nil {
		s.DeleteTarget(te)
		return err
	}
	if te.IsECTarget() {
		ecsvc.GetECService().AssociateTarget(te.Id, te.CustomLabel)
	}
	return nil
}

func isCategoryEmpty(request *dto.TargetModel) bool {
	return request.TagName.MainCategory == "" && request.TagName.SubCategory == "" && request.TagName.DefectDetail == ""
}

func isCategoryDetailMatch(requestCategory string, targetCategory string) bool {
	if requestCategory == "所有" {
		return true
	}

	return requestCategory == targetCategory
}

func isCategoryMatch(request *dto.TargetModel, target *entity.TargetEntity) bool {
	if isCategoryEmpty(request) {
		return true
	}

	if !isCategoryDetailMatch(request.TagName.MainCategory, target.TagName.MainCategory) {
		return false
	}

	if !isCategoryDetailMatch(request.TagName.SubCategory, target.TagName.SubCategory) {
		return false
	}

	return isCategoryDetailMatch(request.TagName.DefectDetail, target.TagName.DefectDetail)
}

func isLanguageMatch(requestLang string, targetLang string) bool {
	if requestLang == "" {
		return true
	}

	return requestLang == targetLang
}

func isInstituteLabelMatch(requestInstituteLabel []string, targetInstituteLabel []string) bool {
	if len(requestInstituteLabel) == 0 {
		return true
	}
	for _, value := range requestInstituteLabel {
		if tools.IsContain(targetInstituteLabel, value) {
			return true
		}
	}

	return false
}

// 过滤工作空间相同
func isWorkspaceMatch(requestWokspace string, targetWokspace string) bool {
	// 查询参数为空时查询所有的工作空间
	if requestWokspace == "" {
		return true
	}
	return requestWokspace == targetWokspace
}

func (s *TargetService) QueryTargetName(request *dto.TargetModel) (bool, error) {
	temp := false
	targets, err := s.FindAll()
	if err != nil {
		return temp, errors.Wrapf(err, "find all targets")
	}
	for _, t := range targets {
		if t.Name == request.Name {
			temp = true
		}
	}
	return temp, nil
}

func (s *TargetService) QueryTargets(request *dto.TargetModel) ([]entity.TargetEntity, error) {
	var targetEntities []entity.TargetEntity
	if request.ID != "" {
		target, err := s.FindTarget(request.ID)
		if err != nil {
			return nil, errors.Wrapf(err, "find specified targets")
		}
		targetEntities = append(targetEntities, *target)
		return targetEntities, nil
	}

	if !request.IsQueryUserValid() {
		return nil, errors.New("user is empty for query all targets")
	}

	targets, err := s.FindAll()
	if err != nil {
		return nil, errors.Wrapf(err, "find all targets")
	}

	sort.Slice(targets, func(i, j int) bool {
		if !targets[i].UpdateTime.Equal(targets[j].UpdateTime) {
			return targets[i].UpdateTime.After(targets[j].UpdateTime)
		}
		return targets[i].CreateTime.After(targets[j].CreateTime)
	})

	for _, t := range targets {
		if isUserCannotSeeTarget := !request.QueryUserIsTargetOwner(t.Owner) && !t.IsShared; isUserCannotSeeTarget {
			continue
		}

		if isWorkspaceMatch(request.Workspace, t.Workspace) && isLanguageMatch(request.Lang, t.Language) &&
			isCategoryMatch(request, &t) && isInstituteLabelMatch(request.InstituteLabel, t.InstituteLabel) {
			targetEntities = append(targetEntities, t)
		}
	}
	return targetEntities, nil
}

func (s *TargetService) RemoveTarget(request *dto.TargetModel) error {
	if request.Owner == "" {
		return errors.New("owner is empty")
	}
	if request.ID == "" {
		return errors.New("target id is empty")
	}

	target, err := s.FindTarget(request.ID)
	if err != nil || target == nil {
		return errors.WithMessagef(err, "find target %v", request.ID)
	}
	if target.Owner != request.Owner {
		return errors.Errorf("removing target is forbidden")
	}
	if len(target.RelatedRanges) != 0 {
		return errors.Errorf("target %v has related ranges", request.ID)
	}
	if err := os.RemoveAll(target.GetTargetDir()); err != nil && !os.IsNotExist(err) {
		return errors.Errorf("removing target dir fails, err: %v", err)
	}
	if target.IsECTarget() {
		ecsvc.GetECService().DisassociateTarget(target.Id, target.CustomLabel)
	}
	return s.DeleteTarget(target)
}

/*
*

	批量导出靶子操作
*/
func (s *TargetService) ExportBatchTarget(targetIds []string) error {
	if len(targetIds) == 0 {
		return errors.New("target id is empty")
	}
	zip_file_name := "/tmp/targets.zip"
	zipfile, _ := os.Create(zip_file_name)
	defer zipfile.Close()
	archive := zip.NewWriter(zipfile)
	defer archive.Close()
	for i := 0; i < len(targetIds); i++ {
		target, err := s.FindTarget(targetIds[i])
		if err != nil || target == nil {
			return errors.WithMessagef(err, "find target %v", targetIds[i])
		}
		dir, err := ioutil.ReadDir(target.GetTargetDir())
		if err != nil {
			return err
		}
		if len(dir) == 0 {
			return nil
		}
		filepath.Walk(target.GetTargetDir(), func(path string, info os.FileInfo, _ error) error {
			if path == target.GetTargetDir() {
				return nil
			}
			header, _ := zip.FileInfoHeader(info)
			if strings.Contains(target.Name, "/") {
				header.Name = "targets/" + strings.ReplaceAll(target.Name, "/", "_") + "/" + path[len(target.GetTargetDir())+1:]
			} else if strings.Contains(target.Name, "&") {
				header.Name = "targets/" + strings.ReplaceAll(target.Name, "&", "_") + "/" + path[len(target.GetTargetDir())+1:]
			} else {
				header.Name = "targets/" + target.Name + "/" + path[len(target.GetTargetDir())+1:]
			}

			if info.IsDir() {
				header.Name += `/`
			} else {
				header.Method = zip.Deflate
			}
			writer, _ := archive.CreateHeader(header)
			if !info.IsDir() {
				file, _ := os.Open(path)
				defer file.Close()
				io.Copy(writer, file)
			}
			return nil
		})
	}
	return nil
}
func (s *TargetService) UploadTmpFile(file multipart.File, fileType, name, ownerId string) error {
	return tools.SaveFile(filepath.Join(tempDir(ownerId, fileType), name), file)
}

func (s *TargetService) delRelatedRange(targetID string, rangeID string) error {
	t, err := s.FindTarget(targetID)
	if err != nil {
		return errors.WithMessagef(err, "find target %v failed", targetID)
	}
	if tools.IsContain(t.RelatedRanges, rangeID) {
		t.RelatedRanges = tools.ListRemoveOne(t.RelatedRanges, rangeID)

		if err := s.updateTargetTotalAnswerInfo(t); err != nil {
			return err
		}

		if err := s.UpdateTarget(t); err != nil {
			return errors.WithMessagef(err, "update target %v failed", targetID)
		}
	}
	return nil
}

func (s *TargetService) updateTargetTotalAnswerInfo(t *entity.TargetEntity) error {
	answerFile := filepath.Join(t.GetAnswerFileDir(), t.Answer)
	var sr *sr.ResultCalculator
	num, score, err := sr.GetAnswerRingNumAndScore(answerFile)
	if err != nil {
		logger.Warnf("Get target nums and score by answerFile failed when delete range: %v", answerFile, err.Error())
		return err
	}
	t.TotalAnswerNum = int(num)
	t.TotalAnswerScore = int(score)
	return nil
}

func (s *TargetService) addRelatedRange(targetID string, rangeID string) error {
	t, err := s.FindTarget(targetID)
	if err != nil {
		return errors.WithMessagef(err, "find target %v failed", targetID)
	}
	if !tools.IsContain(t.RelatedRanges, rangeID) {
		t.RelatedRanges = append(t.RelatedRanges, rangeID)

		if err := s.updateTargetTotalAnswerInfo(t); err != nil {
			return err
		}

		if err := s.UpdateTarget(t); err != nil {
			return errors.WithMessagef(err, "update target %v failed", targetID)
		}
	}
	return nil
}

func (s *TargetService) ModifyRelatedRanges(rangeID string, oldtargetIDs, newtargetIDs []string) error {
	for _, oldtargetId := range oldtargetIDs {
		if !tools.IsContain(newtargetIDs, oldtargetId) {
			if err := s.delRelatedRange(oldtargetId, rangeID); err != nil {
				return err
			}
		}
	}
	for _, newtargetId := range newtargetIDs {
		if !tools.IsContain(oldtargetIDs, newtargetId) {
			if err := s.addRelatedRange(newtargetId, rangeID); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *TargetService) ModifyTarget(tm *dto.TargetModel) error {
	t, err := s.FindTarget(tm.ID)
	if err != nil {
		return errors.WithMessagef(err, "find target %v failed", tm.ID)
	}
	isECBefore, ecIdBefore := t.IsECTarget(), t.CustomLabel
	oldFiles := t.Targets
	t.Name = tm.Name
	t.Template = tm.Template
	t.IsShared = tm.Isshared
	t.TagId = tm.TagID
	t.Targets = tm.Files
	t.Answer = tm.Answer
	t.Language = tm.Lang
	t.UpdateTime = time.Now()
	t.TagName.MainCategory = tm.TagName.MainCategory
	t.TagName.SubCategory = tm.TagName.SubCategory
	t.TagName.DefectDetail = tm.TagName.DefectDetail
	t.CustomLabel = tm.CustomLabel
	t.ExtendedLabel = tm.ExtendedLabel
	t.InstituteLabel = tm.InstituteLabel
	t.Workspace = tm.Workspace

	if err := s.updateTargetTotalAnswerInfo(t); err != nil {
		return err
	}

	if err := s.UpdateTarget(t); err != nil {
		return errors.WithMessagef(err, "update target %v failed", tm.ID)
	}
	s.refreshAssociationWithEC(isECBefore, ecIdBefore, t)
	return t.RemoveUselessCodeFiles(oldFiles)
}

func (s *TargetService) refreshAssociationWithEC(isECBefore bool, ecIdBefore string, curTarget *entity.TargetEntity) {
	if isECBefore && (!curTarget.IsECTarget() || ecIdBefore != curTarget.CustomLabel) {
		ecsvc.GetECService().DisassociateTarget(curTarget.Id, ecIdBefore)
	}
	if curTarget.IsECTarget() && (!isECBefore || ecIdBefore != curTarget.CustomLabel) {
		ecsvc.GetECService().AssociateTarget(curTarget.Id, curTarget.CustomLabel)
	}
}

func (s *TargetService) UpdateFile(id, fileType, name string, content io.Reader) error {
	t, err := s.FindTarget(id)
	if err != nil {
		return errors.WithMessagef(err, "find target %v failed", id)
	}
	if fileType == entity.TypeTarget {
		err = t.UpdateCodeFile(name, content)
	} else {
		err = t.UpdateAnswerFile(name, content)
	}
	if err != nil {
		return err
	}

	t.UpdateTime = time.Now()

	if err := s.updateTargetTotalAnswerInfo(t); err != nil {
		return err
	}

	if err := s.UpdateTarget(t); err != nil {
		return errors.WithMessagef(err, "update target %v failed", id)
	}
	return nil
}

func (s *TargetService) GetCodeFile(id, file string) (io.ReadCloser, error) {
	t, err := s.FindTarget(id)
	if err != nil {
		return nil, errors.WithMessagef(err, "find target %v failed", id)
	}
	if !tools.IsContain(t.Targets, file) {
		return nil, errors.WithMessagef(errcode.ErrRecordNotFound, "target code file %v not found", file)
	}
	f, err := os.Open(filepath.Join(t.GetCodeFilesDir(), file))
	if err != nil {
		return nil, errors.WithMessagef(errcode.ErrFileSystemError, "open code file failed: %v", err)
	}
	return f, nil
}

func (s *TargetService) FindByID(id string) (*entity.TargetEntity, error) {
	return s.FindTarget(id)
}

func (s *TargetService) FindAllTarget() ([]entity.TargetEntity, error) {
	tar := GetTargetService()
	return tar.FindAll()
}

func (s *TargetService) GetTargetAnswers(id string) ([]shootingresult.TargetAnswer, error) {
	t, err := s.FindTarget(id)
	if err != nil {
		return nil, errors.WithMessagef(err, "find target %v failed", id)
	}
	//temp get answers for show to user, got from repository after answer-built system
	return shootingresult.NewShootingResultCalculator().LoadShootingData(filepath.Join(t.GetAnswerFileDir(), t.Answer))
}
