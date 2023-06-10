package ecsvc

import (
	"code-shooting/domain/entity"
	"code-shooting/domain/entity/ec"
	"code-shooting/domain/entity/spec"
	"code-shooting/domain/repository"

	"code-shooting/infra/logger"
)

type ECService struct {
}

func GetECService() *ECService {
	return &ECService{}
}

func (s *ECService) AddEC(e *ec.EC) error {
	ts, err := s.searchTargets(e.Id)
	if err != nil {
		return err
	}
	for _, t := range ts {
		e.AssociateWithTarget(t.Id)
	}
	return ec.GetECRepo().Save(e)
}

func (s *ECService) searchTargets(ecId string) ([]entity.TargetEntity, error) {
	targets, err := repository.GetTargetRepo().Find(spec.NewAndSpec(
		spec.ExtendedLabels.AnyEq(entity.ExtLabelEC), spec.CustomLabel.Equal(ecId),
	))
	if err != nil {
		return nil, err
	}
	return targets, nil
}

func (s *ECService) AssociateTarget(targetId, ecId string) {
	err := s.updateEC(ecId, func(e *ec.EC) {
		e.AssociateWithTarget(targetId)
	})
	if err != nil {
		logger.Warnf("associate target %v with ec %v failed: %v", targetId, ecId, err)
	}
}

func (s *ECService) DisassociateTarget(targetId, ecId string) {
	err := s.updateEC(ecId, func(e *ec.EC) {
		e.DisassociateWithTarget(targetId)
	})
	if err != nil {
		logger.Warnf("associate target %v with ec %v failed: %v", targetId, ecId, err)
	}
}

func (s *ECService) updateEC(ecId string, f func(e *ec.EC)) error {
	e, err := ec.GetECRepo().Get(ecId)
	if err != nil {
		return err
	}
	f(e)
	return ec.GetECRepo().Save(e)
}
