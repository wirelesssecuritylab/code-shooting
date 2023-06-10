package sql

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "github.com/lib/pq" // #nolint
	"github.com/pkg/errors"

	"code-shooting/infra/internal/idatabase"
	"code-shooting/infra/logger"
)

type DB struct {
	sync.Mutex
	*sql.DB
	cfg config
}

func (s *DB) init() error {
	err := s.connect(s.cfg.Password)
	if err != nil {
		logger.Error("failed connect db of ", s.cfg.ID, " : ", err)
		return err
	}
	logger.Info("success init pg of ", s.cfg.ID)
	return nil
}

func (s *DB) IsReady() bool {
	s.Lock()
	defer s.Unlock()
	if s.DB == nil {
		logger.Warn(s.cfg.ID, " is not initialized")
		return false
	}
	err := s.DB.Ping()
	if err != nil {
		logger.Warn(s.cfg.ID, " is not ready: ", err)
		return false
	}
	return true
}

func (s *DB) Close() error {
	if s.DB == nil {
		return nil
	}
	return s.DB.Close()
}

func (s *DB) connect(pwd string) error {
	s.Lock()
	defer s.Unlock()

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s "+
		"dbname=%s",
		s.cfg.Host, s.cfg.Port, s.cfg.User, pwd,
		s.cfg.DBName)
	for k, v := range s.cfg.ConnParams {
		dsn = dsn + " " + k + "=" + v
	}
	db, err := idatabase.OpenSqlDB(dsn)
	if err != nil {
		logger.Error("open db error of id ", s.cfg.ID, " : ", err)
		return errors.Wrapf(err, "open db error of id %s", s.cfg.ID)
	}

	if err = db.Ping(); err != nil {
		logger.Error("connect db error of id ", s.cfg.ID, " : ", err)
		return errors.Wrapf(err, "connect db error of id %s", s.cfg.ID)
	}

	for _, option := range s.cfg.DBOptions {
		option(db)
	}

	s.closeSqlDB(s.DB)

	s.DB = db

	logger.Info("success connect db of id ", s.cfg.ID)
	return nil
}

func (s *DB) updateLink(pwd string) error {
	logger.Info("update es native client of ", s.cfg.ID)
	err := s.connect(pwd)
	if err != nil {
		return errors.Wrapf(err, "update link err ")
	}
	return nil
}

func (s *DB) closeSqlDB(db *sql.DB) {
	if db != nil {
		time.AfterFunc(1*time.Minute, func() {
			err := db.Close()
			logger.Info("close db of ", s.cfg.ID, " : ", err)
		})
	}
}

func (s *DB) GetDBName() string {
	return s.cfg.DBName
}
