package ecsvc

import (
	"encoding/csv"
	"io"
	"reflect"
	"time"

	"github.com/pkg/errors"

	"code-shooting/domain/entity/ec"
	"code-shooting/domain/entity/spec"
	ecsvc "code-shooting/domain/service/ec"
	"code-shooting/infra/errcode"
	"code-shooting/interface/dto"
)

type ECService struct {
}

func GetECService() *ECService {
	return &ECService{}
}

type ECImportInfo struct {
	Institute string `json:"institute"`
	Center    string `json:"center"`
	UserId    string `json:"userId"`
}

func (s *ECService) ImportFromCsv(ei *ECImportInfo, cr *csv.Reader) error {
	if err := s.readTitle(cr); err != nil {
		return err
	}
	return s.readRecords(ei, cr)
}

func (s *ECService) readTitle(cr *csv.Reader) error {
	title, err := cr.Read()
	if err != nil {
		return err
	}
	if !reflect.DeepEqual(title, []string{"EC号", "部门", "团队"}) {
		return errors.WithMessage(errcode.ErrParamError, "invalid ec records title")
	}
	return nil
}

func (s *ECService) readRecords(ei *ECImportInfo, cr *csv.Reader) error {
	for {
		record, err := cr.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return errors.WithMessagef(errcode.ErrParamError, "invalid ec record: %v", err)
		}
		if err := s.importEC(ei, record); err != nil {
			return err
		}
	}
	return nil
}

func (s *ECService) importEC(ei *ECImportInfo, record []string) error {
	return ecsvc.GetECService().AddEC(ec.NewEC(record[0], ec.Organization{
		Institute:  ei.Institute,
		Center:     ei.Center,
		Department: record[1],
		Team:       record[2],
	}, ei.UserId))
}

type TimeFilter struct {
	SinceTime time.Time
	UntilTime time.Time
}

func (s *ECService) QueryECs(org *ec.Organization, tr *TimeFilter) ([]dto.EC, error) {
	as := spec.NewAndSpec(spec.Institute.Equal(org.Institute), spec.Center.Equal(org.Center))
	as.AddIf(org.Department != "", spec.Department.Equal(org.Department))
	as.AddIf(org.Team != "", spec.Department.Equal(org.Team))
	as.AddIf(!tr.SinceTime.IsZero(), spec.ImportTime.Since(tr.SinceTime))
	as.AddIf(!tr.UntilTime.IsZero(), spec.ImportTime.Until(tr.UntilTime))
	ecs, err := ec.GetECRepo().Find(as)
	if err != nil {
		return nil, err
	}
	return s.toDto(ecs), nil
}

func (s *ECService) toDto(ecs []*ec.EC) []dto.EC {
	decs := make([]dto.EC, 0, len(ecs))
	for _, e := range ecs {
		decs = append(decs, dto.EC{
			Id:                e.Id,
			Organization:      e.Organization,
			Importer:          e.Importer,
			ImportTime:        e.ImportTime.Format(time.RFC3339),
			ConvertedToTarget: e.ConvertedToTarget,
		})
	}
	return decs
}

func (s *ECService) DeleteEC(id, userId string) error {
	return ec.GetECRepo().Remove(id, userId)
}
