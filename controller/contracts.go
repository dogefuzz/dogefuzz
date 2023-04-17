package controller

import (
	"net/http"

	"github.com/dogefuzz/dogefuzz/environment"
	"github.com/dogefuzz/dogefuzz/pkg/interfaces"
	"github.com/gin-gonic/gin"
)

type GetAgentsResponse struct {
	Addresses []string `json:"addresses"`
}

type contractsController struct {
	contractService interfaces.ContractService
}

func NewContractsController(e Env) *contractsController {
	return &contractsController{
		contractService: e.ContractService(),
	}
}

func (ctrl *contractsController) GetAgents(c *gin.Context) {
	agentsIds := []string{
		environment.EXCEPTION_AGENT_CONTRACT_ID,
		environment.GAS_CONSUMPTION_AGENT_CONTRACT_ID,
		environment.REENTRANCY_AGENT_CONTRACT_ID,
	}

	addresses := make([]string, len(agentsIds))
	for idx, agentId := range agentsIds {
		contract, err := ctrl.contractService.Get(agentId)
		if err != nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		addresses[idx] = contract.Address
	}

	c.JSON(http.StatusOK, GetAgentsResponse{
		Addresses: addresses,
	})
}
