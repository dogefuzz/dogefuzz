package api

func BuildRoutes(s *server) {
	gr := s.router.Group("/")

	gr.GET("/contracts/agent", s.env.ContractsController().GetAgent)
	gr.POST("/tasks", s.env.TasksController().Start)
	gr.POST("/transactions/executions", s.env.TransactionsController().StoreTransactionExecution)
	gr.POST("/transactions/weaknesses", s.env.TransactionsController().StoreDetectedWeaknesses)
}
