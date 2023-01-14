package it

import (
	"fmt"

	"github.com/dogefuzz/dogefuzz/config"
	"github.com/dogefuzz/dogefuzz/data"
	"github.com/dogefuzz/dogefuzz/pkg/interfaces"
	"go.uber.org/zap"
)

type Env struct {
	cfg          *config.Config
	logger       *zap.Logger
	dbConnection interfaces.Connection
}

func NewEnv(cfg *config.Config) *Env {
	return &Env{cfg: cfg}
}

func (e *Env) Destroy() {
	if e.dbConnection != nil {
		e.dbConnection.Clean()
	}
	if e.dbConnection != nil {
		e.logger.Sync()
	}
}

func (e *Env) Config() *config.Config {
	return e.cfg
}

func (e *Env) Logger() *zap.Logger {
	if e.logger == nil {
		e.logger = zap.NewNop()
	}
	return e.logger
}

func (e *Env) DbConnection() interfaces.Connection {
	if e.dbConnection == nil {
		dbConnection, err := data.NewConnection(e.Config(), e.Logger())
		if err != nil {
			e.logger.Error(fmt.Sprintf("Error while initializing database manager: %s", err))
			return nil
		}
		e.dbConnection = dbConnection
	}
	return e.dbConnection
}
