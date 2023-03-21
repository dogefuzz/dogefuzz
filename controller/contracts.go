package controller

import (
	"net/http"

	"github.com/dogefuzz/dogefuzz/environment"
	"github.com/dogefuzz/dogefuzz/pkg/interfaces"
	"github.com/gin-gonic/gin"
)

type GetAgentResponse struct {
	Address string `json:"address"`
}

type contractsController struct {
	contractService interfaces.ContractService
}

func NewContractsController(e Env) *contractsController {
	return &contractsController{
		contractService: e.ContractService(),
	}
}

func (ctrl *contractsController) GetAgent(c *gin.Context) {
	contract, err := ctrl.contractService.Get(environment.REENTRANCY_AGENT_CONTRACT_ID)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, GetAgentResponse{
		Address: contract.Address,
	})
}
