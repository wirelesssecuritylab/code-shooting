package idatabase

import "database/sql"

func OpenSqlDB(dataSourceName string) (*sql.DB, error) {
	return sql.Open("postgres", dataSourceName)
}
