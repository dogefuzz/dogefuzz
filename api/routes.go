package api

func BuildRoutes(s *server) {
	gr := s.router.Group("/")

	gr.POST("/tasks", s.env.TasksController().Start)
	gr.PUT("/tasks/:task_id", s.env.TasksController().Stop)
	gr.POST("/executions", s.env.ExecutionsController().Create)
	gr.POST("/weaknesses", s.env.WeaknessesController().Create)
	gr.POST("/transactions", s.env.TransactionsController().Create)
}
