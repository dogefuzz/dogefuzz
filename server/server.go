package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gongbell/contractfuzzer/env"
	"github.com/gongbell/contractfuzzer/server/api"
)

type Server interface {
	Run() error
}

type DefaultServer struct {
	environment     env.Environment
	port            string
	ginServer       *http.Server
	fuzzAPIs        api.FuzzAPI
	hackAPIs        api.HackAPI
	instrumentAPIs  api.InstrumentAPI
	transactionAPIs api.TransactionAPI
}

func (s DefaultServer) Init(environment env.Environment, addrMapPath, reporter, port string) DefaultServer {
	s.environment = environment
	s.port = port

	s.fuzzAPIs = new(api.DefaultFuzzAPI).Init(
		environment.Logger(),
		environment.EventBus(),
		environment.TaskRepository(),
		environment.OracleRepository(),
		environment.TaskOracleRepository(),
		environment.ContractRepository(),
		environment.TaskContractRepository(),
	)
	s.instrumentAPIs = new(api.DefaultInstrumentAPI).Init(
		environment.Logger(),
		environment.EventBus(),
		environment.TransactionRepository(),
		environment.TaskRepository(),
		environment.TaskOracleRepository(),
		environment.OracleRepository(),
	)
	s.transactionAPIs = new(api.DefaultTransactionAPI).Init(
		environment.Logger(),
		environment.TransactionRepository(),
		environment.ContractRepository(),
	)

	// Configure router
	router := gin.Default()
	router.POST("/fuzz/start", s.fuzzAPIs.Start)
	router.POST("/fuzz/stop/:task_id", s.fuzzAPIs.Stop)
	router.POST("/instrument/execution", s.instrumentAPIs.Execution)
	router.POST("/instrument/weakness", s.instrumentAPIs.Weakness)
	router.POST("/transaction", s.transactionAPIs.Create)

	// Configure internal server
	s.ginServer = &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: router,
	}
	return s
}

func (s DefaultServer) Run() error {
	s.environment.Logger().Info(fmt.Sprintf("Running server in localhost:%s", s.port))
	return s.ginServer.ListenAndServe()
}

func (s DefaultServer) Shutdown(ctx context.Context) error {
	return s.ginServer.Shutdown(ctx)
}
