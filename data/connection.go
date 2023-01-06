package data

import (
	"fmt"
	"os"

	"github.com/dogefuzz/dogefuzz/config"
	"github.com/dogefuzz/dogefuzz/entities"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Connection interface {
	Migrate() error
	Clean() error
	GetDB() *gorm.DB
}

type connection struct {
	db           *gorm.DB
	databaseName string
	logger       *zap.Logger
}

func NewConnection(cfg *config.Config, logger *zap.Logger) (*connection, error) {
	logger.Info(fmt.Sprintf("Initializing database in \"%s\" file", cfg.DatabaseName))
	db, err := gorm.Open(sqlite.Open(cfg.DatabaseName), &gorm.Config{})
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

func (m *connection) GetDB() *gorm.DB {
	return m.db
}

func (m *connection) Migrate() error {
	return m.db.AutoMigrate(&entities.Contract{}, &entities.Function{}, &entities.Task{}, &entities.Transaction{})
}
