package database

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // justifying it
)

type DBConnector interface {
	Client() *sqlx.DB
}
