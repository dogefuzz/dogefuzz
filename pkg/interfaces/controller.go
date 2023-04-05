package interfaces

import "github.com/gin-gonic/gin"

type ContractsController interface {
	GetAgent(c *gin.Context)
}

type TasksController interface {
	Start(c *gin.Context)
}

type TransactionsController interface {
	StoreDetectedWeaknesses(c *gin.Context)
	StoreTransactionExecution(c *gin.Context)
}

type PingController interface {
	Ping(c *gin.Context)
}
