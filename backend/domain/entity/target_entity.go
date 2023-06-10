package entity

import (
	"io"
	"os"
	"path/filepath"

	"code-shooting/infra/logger"

	"code-shooting/infra/util"
	"code-shooting/infra/util/tools"
	"time"
)

const (
	_TargetsDir = "targets"
	TypeTarget  = "target"
	TypeAnswer  = "answer"

	ExtLabelEC = "外场故障"
)

type TargetEntity struct {
	Id               string      `json:"id"`
	Name             string      `json:"name"`
	Language         string      `json:"language"`
	Template         string      `json:"template"`
	Owner            string      `json:"owner"`
	OwnerName        string      `json:"ownerName"`
	IsShared         bool        `json:"isShared"`
	TagId            string      `json:"tagId,omitempty"`
	TagName          TagNameInfo `json:"tagName,omitempty"`
	CustomLabel      string      `json:"customLable,omitempty"`
	ExtendedLabel    []string    `json:"extendedLabel,omitempty"`
	InstituteLabel   []string    `json:"instituteLabel,omitempty"`
	Answer           string      `json:"answer,omitempty"`
	Targets          []string    `json:"targets,omitempty"`
	RelatedRanges    []string    `json:"relatedRanges,omitempty"`
	Workspace        string      `json:"workspace"`
	CreateTime       time.Time   `json:"createTime,omitempty"`
	UpdateTime       time.Time   `json:"updateTime,omitempty"`
	TotalAnswerNum   int         `json:"totalAnswerNum,omitempty"`
	TotalAnswerScore int         `json:"totalAnswerScore,omitempty"`
}

type TagNameInfo struct {
	MainCategory string `json:"mainCategory,omitempty"`
	SubCategory  string `json:"subCategory,omitempty"`
	DefectDetail string `json:"defectDetail,omitempty"`
}

func (s *TargetEntity) RemoveUselessCodeFiles(oldFiles []string) error {
	for _, f := range oldFiles {
		if tools.IsContain(s.Targets, f) {
			continue
		}
		err := os.Remove(filepath.Join(s.GetCodeFilesDir(), f))
		if err != nil {
			logger.Warnf("remove target %v code file %v failed: %v", s.Id, f, err)
		}
	}
	return nil
}

func (s *TargetEntity) UpdateCodeFile(name string, content io.Reader) error {
	err := tools.SaveFile(filepath.Join(s.GetCodeFilesDir(), name), content)
	if err != nil {
		return err
	}
	if !tools.IsContain(s.Targets, name) {
		s.Targets = append(s.Targets, name)
	}
	return nil
}

func (s *TargetEntity) UpdateAnswerFile(name string, content io.Reader) error {
	err := tools.SaveFile(filepath.Join(s.GetAnswerFileDir(), name), content)
	if err != nil {
		return err
	}
	s.Answer = name
	return nil
}

func (s *TargetEntity) GetCodeFilesDir() string {
	return filepath.Join(s.GetTargetDir(), TypeTarget)
}

func (s *TargetEntity) GetAnswerFileDir() string {
	return filepath.Join(s.GetTargetDir(), TypeAnswer)
}

func (s *TargetEntity) GetTargetDir() string {
	return filepath.Join(util.DataDir, _TargetsDir, s.Id)
}

func (s *TargetEntity) IsECTarget() bool {
	return tools.IsContain(s.ExtendedLabel, ExtLabelEC) && s.CustomLabel != ""
}
