package sql

import (
	"database/sql"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GormDB struct {
	sync.Mutex
	*gorm.DB
}

func newGormDB(db *DB) (*GormDB, error) {
	var gormDB = &GormDB{}

	var err error
	gormDB.DB, err = openGormDB(db.DB)
	if err != nil {
		return nil, err
	}

	return gormDB, nil
}

func (s *GormDB) Close() error {
	if s.DB == nil {
		return nil
	}
	sqlDB, err := s.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func openGormDB(db *sql.DB) (*gorm.DB, error) {
	return gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
}
