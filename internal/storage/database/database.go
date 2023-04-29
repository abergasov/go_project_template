package database

import (
	"fmt"
	"go_project_template/internal/config"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type DBConnect struct {
	db *sqlx.DB
}

func InitDBConnect(cnf *config.DBConf, migratesFolder string) (*DBConnect, error) {
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
	conn := &DBConnect{db}
	if migratesFolder != "" {
		if err = conn.migrate(migratesFolder); err != nil && err != migrate.ErrNoChange {
			return nil, fmt.Errorf("error migrate db: %w", err)
		}
	}
	return conn, nil
}

func InitSQLiteDBConnect(dbPath string) (DBConnector, error) {
	db, err := sqlx.Connect("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("error connect to db: %w", err)
	}
	return &DBConnect{db}, err
}

func InitSQLiteDBConnectMemory() (DBConnector, error) {
	db, err := sqlx.Connect("sqlite3", ":memory:")
	if err != nil {
		return nil, fmt.Errorf("error connect to db: %w", err)
	}
	return &DBConnect{db}, err
}

func (d *DBConnect) Close() error {
	return d.db.Close()
}

func (d *DBConnect) Client() *sqlx.DB {
	return d.db
}

func (d *DBConnect) migrate(migratesFolder string) error {
	driver, err := postgres.WithInstance(d.db.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("error generate driver for db migrator: %w", err)
	}
	m, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file://%s", migratesFolder), "postgres", driver)
	if err != nil {
		return fmt.Errorf("error init db migrator: %w", err)
	}
	return m.Up()
}
