package interfaces

import "github.com/gin-gonic/gin"

type TasksController interface {
	Start(c *gin.Context)
}

type TransactionsController interface {
	StoreDetectedWeaknesses(c *gin.Context)
	StoreTransactionExecution(c *gin.Context)
}
