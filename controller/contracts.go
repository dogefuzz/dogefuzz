package controller

import (
	"net/http"

	"github.com/dogefuzz/dogefuzz/dto"
	"github.com/dogefuzz/dogefuzz/pkg/solc"
	"github.com/dogefuzz/dogefuzz/service"
	"github.com/gin-gonic/gin"
)

type ContractsController interface {
	Create(c *gin.Context)
}

type contractsController struct {
	contractService  service.ContractService
	solidityCompiler solc.SolidityCompiler
}

func NewContractsController(e Env) *contractsController {
	return &contractsController{
		contractService:  e.ContractService(),
		solidityCompiler: e.SolidityCompiler(),
	}
}

func (ctrl *contractsController) Create(c *gin.Context) {
	var request dto.NewContractDTO
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	contract, err := ctrl.solidityCompiler.CompileSource(request.Source)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
	}
	request.AbiDefinition = contract.AbiDefinition
	request.CompiledCode = contract.CompiledCode
	request.Name = contract.Name

	created, err := ctrl.contractService.Create(&request)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, created)
}
