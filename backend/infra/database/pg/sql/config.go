package sql

import (
	"database/sql"

)

type config struct {
	ID         string
	User       string
	Password   string
	Host       string
	Port       int
	DBName     string
	ConnParams map[string]string
	DBOptions  []func(db *sql.DB)
}
