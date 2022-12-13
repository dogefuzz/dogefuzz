package db

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/dogefuzz/dogefuzz/config"
	"go.uber.org/zap"
)

type Connection interface {
	Migrate() error
	Clean() error
	GetDB() *sql.DB
}

type connection struct {
	db           *sql.DB
	databaseName string
	logger       *zap.Logger
}

func NewConnection(cfg *config.Config, logger *zap.Logger) (*connection, error) {
	logger.Info(fmt.Sprintf("Initializing database in \"%s\" file", cfg.DatabaseName))
	db, err := sql.Open("sqlite3", cfg.DatabaseName)
	if err != nil {
		return nil, err
	}

	return &connection{
		db:           db,
		databaseName: cfg.DatabaseName,
		logger:       logger,
	}, nil
}

func (m *connection) Clean() error {
	m.logger.Info(fmt.Sprintf("Cleaning \"%s\" database file", m.databaseName))
	return os.Remove(m.databaseName)
}

func (m *connection) GetDB() *sql.DB {
	return m.db
}

func (m *connection) Migrate() error {
	_, err := m.db.Exec(MIGRATION_QUERY)
	return err
}
