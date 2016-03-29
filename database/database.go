package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

// TODO: load login-info from configuration-file
func NewConnection() (*sql.DB, error) {
	db, err := sql.Open("mysql", "master:badpassword@/cluster")
	if err != nil {
		return nil, err
	}

	return db, nil
}
