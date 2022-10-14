package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dogefuzz/dogefuzz/config"
	"github.com/gin-gonic/gin"
)

type Server interface {
	Start() error
	Shutdown(ctx context.Context) error
}

type server struct {
	server *http.Server
	router *gin.Engine
	cfg    *config.Config
	env    *env
}

func NewServer(cfg *config.Config) *server {
	return &server{
		cfg:    cfg,
		env:    NewEnv(cfg),
		router: gin.Default(),
	}
}

func (s *server) Start() error {
	BuildRoutes(s)
	s.server = &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", s.cfg.ServerPort),
		Handler: s.router,
	}
	s.env.Logger().Info(fmt.Sprintf("Running server in 0.0.0.0:%d", s.cfg.ServerPort))

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
