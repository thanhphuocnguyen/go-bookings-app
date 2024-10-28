package driver

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DB struct {
	SQL *sql.DB
}

var dbConn = &DB{}

const (
	maxOpenDbConn = 10
	maxIdleDbConn = 5
	maxDbLifeTime = 5 * time.Minute
)

func GetDB(connectionString string) (*DB, error) {
	d, err := ConnectDatabase(connectionString)
	if err != nil {
		panic(err)
	}

	d.SetMaxIdleConns(maxIdleDbConn)
	d.SetMaxOpenConns(maxOpenDbConn)
	d.SetConnMaxLifetime(maxDbLifeTime)

	dbConn.SQL = d

	return dbConn, nil
}

func ConnectDatabase(connectionString string) (*sql.DB, error) {
	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
