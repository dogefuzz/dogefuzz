package db

import (
	"database/sql"
	"fmt"
	"os"

	"go.uber.org/zap"
)

type Manager interface {
	Migrate() error
	Seed() error
	Clean() error
	GetDB() *sql.DB
}

type SQLiteManager struct {
	db           *sql.DB
	databaseName string
	logger       *zap.Logger
}

func (m SQLiteManager) Init(logger *zap.Logger) (SQLiteManager, error) {
	m.logger = logger
	m.databaseName = "fuzzer.db"

	m.logger.Info(fmt.Sprintf("Initializing database in \"%s\" file", m.databaseName))
	db, err := sql.Open("sqlite3", m.databaseName)
	if err != nil {
		return SQLiteManager{}, err
	}
	m.db = db

	// run migration query
	err = m.Migrate()
	if err != nil {
		return SQLiteManager{}, err
	}

	// run seed query
	err = m.Seed()
	if err != nil {
		return SQLiteManager{}, err
	}

	return m, nil
}

func (m SQLiteManager) Clean() error {
	m.logger.Info(fmt.Sprintf("Cleaning \"%s\" database file", m.databaseName))
	return os.Remove(m.databaseName)
}

func (m SQLiteManager) GetDB() *sql.DB {
	return m.db
}
