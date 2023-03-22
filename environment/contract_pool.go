package environment

import (
	"context"
	"fmt"
	"io/ioutil"
	"sync"

	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/dogefuzz/dogefuzz/pkg/dto"
	"github.com/dogefuzz/dogefuzz/pkg/interfaces"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

const CONTRACT_FOLDER = "./assets/contracts"

type contractPool struct {
	mu                 sync.Mutex
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
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return p.deployExceptionFallbackContract(ctx)
	})

	g.Go(func() error {
		return p.deployGasConsumptionFallbackContract(ctx)
	})

	g.Go(func() error {
		return p.deployReentrancyAgentContract(ctx)
	})

	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}

func (p *contractPool) deployExceptionFallbackContract(ctx context.Context) error {
	name := "ExceptionFallback"
	return p.deployContract(ctx, name, EXCEPTION_FALLBACK_CONTRACT_ID)
}

func (p *contractPool) deployGasConsumptionFallbackContract(ctx context.Context) error {
	name := "GasConsumptionFallback"
	return p.deployContract(ctx, name, GAS_CONSUMPTION_FALLBACK_CONTRACT_ID)
}

func (p *contractPool) deployReentrancyAgentContract(ctx context.Context) error {
	name := "ReentrancyAgent"
	return p.deployContract(ctx, name, REENTRANCY_AGENT_CONTRACT_ID)
}

func (p *contractPool) deployContract(ctx context.Context, contractName string, contractId string) error {
	content, err := ioutil.ReadFile(fmt.Sprintf("%s/%s.sol", CONTRACT_FOLDER, contractName))
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

	p.mu.Lock()
	address, _, err := p.gethService.Deploy(ctx, compiledContract)
	if err != nil {
		return err
	}
	p.mu.Unlock()
	p.logger.Sugar().Debugf("deploying contract %s at %s", contractName, address)

	contract.Address = address
	contract.Status = common.CONTRACT_DEPLOYED
	err = p.contractService.Update(contract)
	if err != nil {
		return err
	}
	return nil
}
