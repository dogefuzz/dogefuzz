package environment

import (
	"context"
	"fmt"
	"os"

	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/dogefuzz/dogefuzz/pkg/dto"
	"github.com/dogefuzz/dogefuzz/pkg/interfaces"
	"go.uber.org/zap"
)

const CONTRACT_FOLDER = "./assets/contracts"

type contractPool struct {
	logger             *zap.Logger
	gethService        interfaces.GethService
	contractService    interfaces.ContractService
	transactionService interfaces.TransactionService
	solidityCompiler   interfaces.SolidityCompiler
}

func NewContractPool(e env) *contractPool {
	return &contractPool{
		logger:             e.Logger(),
		gethService:        e.GethService(),
		contractService:    e.ContractService(),
		transactionService: e.TransactionService(),
		solidityCompiler:   e.SolidityCompiler(),
	}
}

func (p *contractPool) Setup(ctx context.Context) error {
	contracts := map[string]string{
		EXCEPTION_FALLBACK_CONTRACT_NAME:       EXCEPTION_FALLBACK_CONTRACT_ID,
		EXCEPTION_AGENT_CONTRACT_NAME:          EXCEPTION_AGENT_CONTRACT_ID,
		GAS_CONSUMPTION_FALLBACK_CONTRACT_NAME: GAS_CONSUMPTION_FALLBACK_CONTRACT_ID,
		GAS_CONSUMPTION_AGENT_CONTRACT_NAME:    GAS_CONSUMPTION_AGENT_CONTRACT_ID,
		REENTRANCY_AGENT_CONTRACT_NAME:         REENTRANCY_AGENT_CONTRACT_ID,
	}

	for name, id := range contracts {
		err := p.deployContract(ctx, name, id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *contractPool) deployContract(ctx context.Context, contractName string, contractId string) error {
	content, err := os.ReadFile(fmt.Sprintf("%s/%s.sol", CONTRACT_FOLDER, contractName))
	if err != nil {
		return err
	}
	contractSource := string(content)

	compiledContract, err := p.solidityCompiler.CompileSource(contractName, contractSource)
	if err != nil {
		return err
	}

	contract, err := p.contractService.Find(contractId)
	if err != nil {
		return err
	}

	if contract == nil {
		contractDTO := dto.NewContractWithIdDTO{
			Id:                 contractId,
			TaskId:             "",
			Status:             common.CONTRACT_CREATED,
			Source:             contractSource,
			DeploymentBytecode: compiledContract.DeploymentBytecode,
			RuntimeBytecode:    compiledContract.RuntimeBytecode,
			AbiDefinition:      compiledContract.AbiDefinition,
			Name:               compiledContract.Name,
		}
		createdContract, err := p.contractService.CreateWithId(&contractDTO)
		if err != nil {
			return err
		}
		contract = createdContract
	}

	address, _, err := p.gethService.Deploy(ctx, compiledContract)
	if err != nil {
		return err
	}
	p.logger.Sugar().Debugf("deploying contract %s at %s", contractName, address)

	contract.Address = address
	contract.Status = common.CONTRACT_DEPLOYED
	err = p.contractService.Update(contract)
	if err != nil {
		return err
	}
	return nil
}
