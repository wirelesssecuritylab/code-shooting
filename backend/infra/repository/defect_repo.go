package repository

import (
	"code-shooting/infra/po"
	"code-shooting/infra/util/database"
)

type DefectRepository struct {
	DefectDB po.DefectDB
}

func NewDefectRepository() *DefectRepository {
	return &DefectRepository{
		DefectDB: po.DefectDB{GormDB: database.DB},
	}
}

func (s *DefectRepository) SaveInBatch(e *[]po.DefectPo) error {
	return s.DefectDB.SaveDefectInBatch(e)
}

func (s *DefectRepository) RemoveAll() error {
	return s.DefectDB.DeleteAll()
}
