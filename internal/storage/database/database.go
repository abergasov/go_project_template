package database

import (
	"fmt"
	"go_project_template/internal/config"
	"time"

	"github.com/jmoiron/sqlx"
)

type DBConnect struct {
	db *sqlx.DB
}

func InitDBConnect(cnf *config.DBConf) (*DBConnect, error) {
	dsnStr := fmt.Sprintf("dbname=%s sslmode=disable user=%s password=%s host=%s port=%s", cnf.DBName, cnf.User, cnf.Pass, cnf.Address, cnf.Port)
	db, err := sqlx.Connect("postgres", dsnStr)
	if err != nil {
		return nil, fmt.Errorf("error connect to db: %w", err)
	}

	// set db pool max connections
	if cnf.MaxConnections == 0 {
		db.SetMaxOpenConns(10)
	} else {
		db.SetMaxOpenConns(cnf.MaxConnections)
	}
	db.SetConnMaxLifetime(time.Minute)

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error ping to db: %w", err)
	}
	return &DBConnect{db}, nil
}

func (d *DBConnect) Client() *sqlx.DB {
	return d.db
}
