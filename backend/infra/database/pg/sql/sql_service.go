package sql

import (
	"sync"

	"code-shooting/infra/logger"

	marsconfig "code-shooting/infra/config"

	"github.com/pkg/errors"
)

type SqlService interface {
	GetDB(ID string) (*DB, error)
	GetGormDB(ID string) (*GormDB, error)
	CheckStatus() error
	Close() error
	GetAllIds() []string
}

func NewSqlService(config marsconfig.Config) (SqlService, error) {
	cfgs, err := getConfigParsingService().Parse(config)
	if err != nil {
		return nil, errors.Wrap(err, "config content")
	}
	if len(cfgs) == 0 {
		return nil, errors.New("no pg")
	}
	service := &sqlServiceImpl{
		dbs:  make(map[string]*DB),
		gdbs: make(map[string]*GormDB),
	}
	for _, cfg := range cfgs {
		db := &DB{
			cfg: cfg,
		}
		service.dbs[cfg.ID] = db
	}

	return service, nil
}

type sqlServiceImpl struct {
	sync.Mutex
	dbs  map[string]*DB
	gdbs map[string]*GormDB
}

func (s *sqlServiceImpl) CheckStatus() error {
	for _, ID := range s.GetAllIds() {
		_, err := s.GetDB(ID)
		if err != nil {
			logger.Warn("pg is not ready of ", ID, " : ", err)
			return errors.Wrap(err, "check sql service status")
		}
	}
	return nil
}

type errIdMsg struct {
	id  string
	msg string
}

func (s *sqlServiceImpl) Close() error {
	s.Lock()
	defer s.Unlock()
	errIdMsgs := []errIdMsg{}
	for id, db := range s.dbs {
		err := db.Close()
		if err != nil {
			errIdMsgs = append(errIdMsgs, errIdMsg{id: id, msg: err.Error()})
		}
	}
	for id, gdb := range s.gdbs {
		err := gdb.Close()
		if err != nil {
			errIdMsgs = append(errIdMsgs, errIdMsg{id: id, msg: err.Error()})
		}
	}

	if len(errIdMsgs) == 0 {
		return nil
	}

	var errStr string
	for _, errMsg := range errIdMsgs {
		errStr = errStr + errMsg.msg + "of " + errMsg.id
		errStr += "\n"
	}
	return errors.New("db close error: " + errStr)
}

func (s *sqlServiceImpl) GetAllIds() []string {
	s.Lock()
	defer s.Unlock()
	ids := []string{}
	for _, db := range s.dbs {
		ids = append(ids, db.cfg.ID)
	}
	return ids
}

func (s *sqlServiceImpl) GetDB(ID string) (*DB, error) {
	s.Lock()
	defer s.Unlock()
	return s.getDB(ID)
}

func (s *sqlServiceImpl) getDB(ID string) (*DB, error) {
	if db, ok := s.dbs[ID]; ok && db != nil {
		if !db.IsReady() {
			err := db.init()
			if err != nil {
				return nil, errors.Wrapf(err, "init db of %s", ID)
			}

			if !db.IsReady() {
				return nil, errors.Wrapf(err, "db is not ready of %s", ID)
			}
		}
		return db, nil
	}
	return nil, errors.Errorf("%s is not existent", ID)
}

func (s *sqlServiceImpl) GetGormDB(ID string) (*GormDB, error) {
	s.Lock()
	defer s.Unlock()
	if gdb, ok := s.gdbs[ID]; ok && gdb != nil {
		return gdb, nil
	}
	db, err := s.getDB(ID)
	if err != nil {
		return nil, errors.Wrapf(err, "get db of %s", ID)
	}
	var gdb *GormDB
	gdb, err = newGormDB(db)
	if err != nil {
		return nil, errors.Wrapf(err, "init gdb of %s", ID)
	}
	s.gdbs[ID] = gdb
	return gdb, nil
}
