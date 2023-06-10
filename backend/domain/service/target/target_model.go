package target

import (
	"os"
	"path/filepath"

	"code-shooting/infra/logger"
	"github.com/google/uuid"

	"code-shooting/domain/entity"
	"code-shooting/infra/util"
	"code-shooting/interface/dto"
	"time"
)

func NewTarget(t *dto.TargetModel) *entity.TargetEntity {
	return &entity.TargetEntity{
		Id:        uuid.NewString(),
		Name:      t.Name,
		Language:  t.Lang,
		Template:  t.Template,
		Owner:     t.Owner,
		OwnerName: t.OwnerName,
		IsShared:  t.Isshared,
		Targets:   t.Files,
		TagId:     t.TagID,
		Answer:    t.Answer,
		TagName: entity.TagNameInfo{
			MainCategory: t.TagName.MainCategory,
			SubCategory:  t.TagName.SubCategory,
			DefectDetail: t.TagName.DefectDetail,
		},
		CustomLabel:    t.CustomLabel,
		ExtendedLabel:  t.ExtendedLabel,
		InstituteLabel: t.InstituteLabel,
		Workspace:      t.Workspace,
		CreateTime:     time.Now(),
		UpdateTime:     time.Now(),
	}
}

func MoveTmpFiles2Target(t *entity.TargetEntity) error {
	if err := os.MkdirAll(t.GetTargetDir(), 0750); err != nil {
		logger.Warnf("%s", err.Error())
	}
	if err := moveDir(tempDir(t.Owner, entity.TypeTarget), t.GetCodeFilesDir()); err != nil {
		return err
	}
	if err := moveDir(tempDir(t.Owner, entity.TypeAnswer), t.GetAnswerFileDir()); err != nil {
		return err
	}
	return nil
}

func moveDir(srcdir, dstdir string) error {
	return os.Rename(srcdir, dstdir)
}

func tempDir(owner, fileType string) string {
	return filepath.Join(util.DataDir, "temp", owner, fileType)
}
