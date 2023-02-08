package api

import (
	"time"

	ginzap "github.com/gin-contrib/zap"
)

func BuildMiddlewares(s *server) {
	s.router.Use(ginzap.Ginzap(s.env.Logger(), time.RFC3339, true))
	s.router.Use(ginzap.RecoveryWithZap(s.env.Logger(), true))
}
