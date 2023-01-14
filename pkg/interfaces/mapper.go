package interfaces

import (
	"github.com/dogefuzz/dogefuzz/entities"
	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/dogefuzz/dogefuzz/pkg/dto"
)

type ContractMapper interface {
	MapNewDTOToEntity(d *dto.NewContractDTO) *entities.Contract
	MapDTOToEntity(d *dto.ContractDTO) *entities.Contract
	MapEntityToDTO(c *entities.Contract) *dto.ContractDTO
	MapDTOToCommon(c *dto.ContractDTO) *common.Contract
}

type FunctionMapper interface {
	MapNewDTOToEntity(c *dto.NewFunctionDTO) *entities.Function
	MapDTOToEntity(c *dto.FunctionDTO) *entities.Function
	MapEntityToDTO(c *entities.Function) *dto.FunctionDTO
}

type TaskMapper interface {
	MapNewDTOToEntity(c *dto.NewTaskDTO) *entities.Task
	MapDTOToEntity(c *dto.TaskDTO) *entities.Task
	MapEntityToDTO(c *entities.Task) *dto.TaskDTO
}

type TransactionMapper interface {
	MapNewDTOToEntity(n *dto.NewTransactionDTO) *entities.Transaction
	MapDTOToEntity(c *dto.TransactionDTO) *entities.Transaction
	MapEntityToDTO(c *entities.Transaction) *dto.TransactionDTO
}
