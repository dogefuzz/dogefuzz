package controller

import (
	"net/http"

	"github.com/dogefuzz/dogefuzz/dto"
	"github.com/dogefuzz/dogefuzz/service"
	"github.com/gin-gonic/gin"
)

type ContractsController interface {
	Create(c *gin.Context)
}

type contractsController struct {
	contractService service.ContractService
}

func NewContractsController(e Env) *contractsController {
	return &contractsController{
		contractService: e.ContractService(),
	}
}

func (ctrl *contractsController) Create(c *gin.Context) {
	var request dto.NewContractDTO
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	created, err := ctrl.contractService.Create(&request)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, created)
}
