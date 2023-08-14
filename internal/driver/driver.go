package driver

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/lib/pq"
)

// DB holds the database connection pool
type DB struct {
	SQL *sql.DB
}

var dbConn = &DB{}

const maxOpenDbConn = 10
const maxIdleDbConn = 5
const maxDbLifetime = 5 * time.Minute

func ConnectSql(dsn string) (*DB, error) {
	d, err := NewDatabase(dsn)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	d.SetMaxOpenConns(maxOpenDbConn)
	d.SetMaxIdleConns(maxIdleDbConn)
	d.SetConnMaxLifetime(maxDbLifetime)

	dbConn.SQL = d

	// err = testDb(d)

	// if err != nil {
	// 	return nil, err
	// }

	return dbConn, nil
}

func testDb(d *sql.DB) error {
	err := d.Ping()

	if err != nil {
		return err
	}

	return nil
}

func NewDatabase(dsn string) (*sql.DB, error) {
	conn, err := sql.Open("postgres", dsn)

	if err != nil {
		return nil, err
	}

	log.Println("Connected to datastore: ")

	// if err = conn.Ping(); err != nil {
	// 	return nil, err
	// }

	return conn, nil
}
