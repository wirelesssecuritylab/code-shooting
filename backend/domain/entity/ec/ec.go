package ec

import (
	"time"

	"code-shooting/infra/util/tools"
)

type Organization struct {
	Institute  string `json:"institute"`
	Center     string `json:"center"`
	Department string `json:"department"`
	Team       string `json:"team"`
}

type EC struct {
	Organization

	Id                string
	Importer          string
	ImportTime        time.Time
	ConvertedToTarget bool
	AssociatedTargets []string
}

func NewEC(id string, org Organization, importer string) *EC {
	return &EC{
		Id:                id,
		Organization:      org,
		Importer:          importer,
		ImportTime:        time.Now(),
		AssociatedTargets: []string{},
	}
}

func (s *EC) DisassociateWithTarget(target string) {
	s.AssociatedTargets = tools.ListRemoveOne(s.AssociatedTargets, target)
	if len(s.AssociatedTargets) == 0 {
		s.ConvertedToTarget = false
	}
}

func (s *EC) AssociateWithTarget(target string) {
	if !tools.IsContain(s.AssociatedTargets, target) {
		s.AssociatedTargets = append(s.AssociatedTargets, target)
	}
	if !s.ConvertedToTarget {
		s.ConvertedToTarget = true
	}
}
