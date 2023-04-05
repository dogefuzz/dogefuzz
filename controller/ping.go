package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type pingController struct {
}

func NewPingController() *pingController {
	return &pingController{}
}

func (ctrl *pingController) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}
