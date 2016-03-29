package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type DatabaseInfo struct {
	Database         string
	DatabaseUser     string
	DatabasePassword string
}

var info *DatabaseInfo

func Init(db string, usr string, pw string) {
	info = &DatabaseInfo{
		Database:         db,
		DatabaseUser:     usr,
		DatabasePassword: pw,
	}
}

func getConnectionString() string {
	return fmt.Sprintf("%s:%s@/%s", info.DatabaseUser, info.DatabasePassword, info.Database)
}

func NewConnection() (*sql.DB, error) {
	db, err := sql.Open("mysql", getConnectionString())
	if err != nil {
		return nil, err
	}

	return db, nil
}
