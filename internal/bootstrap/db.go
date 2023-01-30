package bootstrap

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"

	"github.com/gocraft/dbr"
)

const (
	maxOpenConns      = 10
	reconnectAttempts = 3
	reconnectInterval = 3 * time.Second
)

// NewDBConnPg creates a new postgres database connection instance.
func NewDBConnPg(username, password, name, host, port string) (*dbr.Connection, error) {
	dsn := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		username,
		password,
		name,
		host,
		port,
	)

	p, err := dbr.Open("postgres", dsn, nil)
	if err != nil {
		return nil, err
	}

	for c := reconnectAttempts; c > 0; c-- {
		err = p.Ping()
		if err != nil {
			time.Sleep(reconnectInterval)

			continue
		}

		break
	}

	if err != nil {
		return nil, err
	}

	p.SetMaxOpenConns(maxOpenConns)

	return p, nil
}

// UpMigrationsPg processes all the pending migrations to the postgres database.
func UpMigrationsPg(db *sql.DB, dbName, path string) (bool, error) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return false, err
	}

	m, err := migrate.NewWithDatabaseInstance("file://"+path, dbName, driver)
	if err != nil {
		return false, err
	}

	err = m.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}
