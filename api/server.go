package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type server struct {
	server *http.Server
	router *gin.Engine
	env    Env
}

func NewServer(env Env) *server {
	// gin.SetMode(gin.ReleaseMode)
	return &server{
		env:    env,
		router: gin.New(),
	}
}

func (s *server) Start() error {
	BuildMiddlewares(s)
	BuildRoutes(s)
	s.server = &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", s.env.Config().ServerPort),
		Handler: s.router,
	}
	s.env.Logger().Info(fmt.Sprintf("running server in 0.0.0.0:%d", s.env.Config().ServerPort))

	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.env.Logger().Sugar().Fatalf("listen: %s\n", err)
		}
	}()

	return nil
}

func (s *server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
