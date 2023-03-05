package generators

import (
	"math/rand"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/dogefuzz/dogefuzz/pkg/dto"
)

func NewContractDTOGen() *dto.NewContractDTO {
	return &dto.NewContractDTO{
		Name:               gofakeit.Name(),
		Source:             gofakeit.LetterN(255),
		DeploymentBytecode: gofakeit.LetterN(255),
		RuntimeBytecode:    gofakeit.LetterN(255),
		AbiDefinition:      gofakeit.LetterN(255),
		TaskId:             gofakeit.LetterN(255),
		Status:             common.RandomChoice([]common.ContractStatus{common.CONTRACT_CREATED, common.CONTRACT_DEPLOYED}),
	}

}

func ContractDTOGen() *dto.ContractDTO {
	return &dto.ContractDTO{
		Id:                 gofakeit.LetterN(255),
		TaskId:             gofakeit.LetterN(255),
		Status:             common.RandomChoice([]common.ContractStatus{common.CONTRACT_CREATED, common.CONTRACT_DEPLOYED}),
		Address:            SmartContractGen(),
		Source:             gofakeit.LetterN(255),
		DeploymentBytecode: gofakeit.LetterN(255),
		AbiDefinition:      gofakeit.LetterN(255),
		Name:               gofakeit.Name(),
		CFG:                CFGGen(),
		DistanceMap:        DistanceMapGen(),
	}
}

func NewFunctionDTOGen() *dto.NewFunctionDTO {
	return &dto.NewFunctionDTO{
		Name:                    gofakeit.LetterN(255),
		NumberOfArgs:            gofakeit.Int64(),
		IsChangingContractState: gofakeit.Bool(),
		IsConstructor:           gofakeit.Bool(),
		ContractId:              gofakeit.LetterN(255),
	}
}

func FunctionDTOGen() *dto.FunctionDTO {
	return &dto.FunctionDTO{
		Id:                      gofakeit.LetterN(255),
		Name:                    gofakeit.LetterN(255),
		NumberOfArgs:            gofakeit.Int64(),
		IsChangingContractState: gofakeit.Bool(),
		IsConstructor:           gofakeit.Bool(),
		ContractId:              gofakeit.LetterN(255),
	}
}

func NewTaskDTOGen() *dto.NewTaskDTO {
	rand.Seed(time.Now().UnixNano())

	argsCount := rand.Intn(255)
	args := make([]string, argsCount)
	for idx := 0; idx < argsCount; idx++ {
		args[idx] = gofakeit.LetterN(255)
	}

	detectorsCount := rand.Intn(255)
	detectors := make([]common.OracleType, detectorsCount)
	for idx := 0; idx < detectorsCount; idx++ {
		detectors[idx] = common.RandomChoice([]common.OracleType{common.DELEGATE_ORACLE, common.EXCEPTION_DISORDER_ORACLE, common.GASLESS_SEND_ORACLE, common.NUMBER_DEPENDENCY_ORACLE, common.REENTRANCY_ORACLE, common.TIMESTAMP_DEPENDENCY_ORACLE})
	}

	aggregatedExecutedInstructionsCount := rand.Intn(255)
	aggregatedExecutedInstructions := make([]string, aggregatedExecutedInstructionsCount)
	for idx := 0; idx < aggregatedExecutedInstructionsCount; idx++ {
		aggregatedExecutedInstructions[idx] = gofakeit.HexUint64()
	}

	return &dto.NewTaskDTO{
		Arguments:                      args,
		StartTime:                      gofakeit.Date(),
		Expiration:                     gofakeit.Date(),
		Detectors:                      detectors,
		FuzzingType:                    common.RandomChoice([]common.FuzzingType{common.BLACKBOX_FUZZING, common.GREYBOX_FUZZING, common.DIRECTED_GREYBOX_FUZZING}),
		AggregatedExecutedInstructions: aggregatedExecutedInstructions,
		Status:                         common.RandomChoice([]common.TaskStatus{common.TASK_DONE, common.TASK_RUNNING}),
	}
}

func TaskDTOGen() *dto.TaskDTO {
	rand.Seed(time.Now().UnixNano())

	argsCount := rand.Intn(255)
	args := make([]string, argsCount)
	for idx := 0; idx < argsCount; idx++ {
		args[idx] = gofakeit.LetterN(255)
	}

	detectorsCount := rand.Intn(255)
	detectors := make([]common.OracleType, detectorsCount)
	for idx := 0; idx < detectorsCount; idx++ {
		detectors[idx] = common.RandomChoice([]common.OracleType{common.DELEGATE_ORACLE, common.EXCEPTION_DISORDER_ORACLE, common.GASLESS_SEND_ORACLE, common.NUMBER_DEPENDENCY_ORACLE, common.REENTRANCY_ORACLE, common.TIMESTAMP_DEPENDENCY_ORACLE})
	}

	aggregatedExecutedInstructionsCount := rand.Intn(255)
	aggregatedExecutedInstructions := make([]string, aggregatedExecutedInstructionsCount)
	for idx := 0; idx < aggregatedExecutedInstructionsCount; idx++ {
		aggregatedExecutedInstructions[idx] = gofakeit.HexUint64()
	}

	return &dto.TaskDTO{
		Id:                             gofakeit.LetterN(255),
		Arguments:                      args,
		StartTime:                      gofakeit.Date(),
		Expiration:                     gofakeit.Date(),
		Detectors:                      detectors,
		FuzzingType:                    common.RandomChoice([]common.FuzzingType{common.BLACKBOX_FUZZING, common.GREYBOX_FUZZING, common.DIRECTED_GREYBOX_FUZZING}),
		AggregatedExecutedInstructions: aggregatedExecutedInstructions,
		Status:                         common.RandomChoice([]common.TaskStatus{common.TASK_DONE, common.TASK_RUNNING}),
	}
}

func NewTransactionDTOGen() *dto.NewTransactionDTO {
	rand.Seed(time.Now().UnixNano())

	inputsCount := rand.Intn(255)
	inputs := make([]string, inputsCount)
	for idx := 0; idx < inputsCount; idx++ {
		inputs[idx] = gofakeit.LetterN(255)
	}

	return &dto.NewTransactionDTO{
		Timestamp:  gofakeit.Date(),
		TaskId:     gofakeit.LetterN(255),
		FunctionId: gofakeit.LetterN(255),
		Inputs:     inputs,
		Status:     common.RandomChoice([]common.TransactionStatus{common.TRANSACTION_DONE, common.TRANSACTION_RUNNING, common.TRANSACTION_CREATED}),
	}
}

func TransactionDTOGen() *dto.TransactionDTO {
	rand.Seed(time.Now().UnixNano())

	inputsCount := rand.Intn(255)
	inputs := make([]string, inputsCount)
	for idx := 0; idx < inputsCount; idx++ {
		inputs[idx] = gofakeit.LetterN(255)
	}

	detectedWeaknessesCount := rand.Intn(255)
	detectedWeaknesses := make([]string, detectedWeaknessesCount)
	for idx := 0; idx < detectedWeaknessesCount; idx++ {
		detectedWeaknesses[idx] = gofakeit.LetterN(255)
	}

	executedInstructionsCount := rand.Intn(255)
	executedInstructions := make([]string, executedInstructionsCount)
	for idx := 0; idx < executedInstructionsCount; idx++ {
		executedInstructions[idx] = gofakeit.LetterN(255)
	}

	return &dto.TransactionDTO{
		Id:                   gofakeit.LetterN(255),
		Timestamp:            gofakeit.Date(),
		BlockchainHash:       gofakeit.LetterN(255),
		TaskId:               gofakeit.LetterN(255),
		FunctionId:           gofakeit.LetterN(255),
		Inputs:               inputs,
		DetectedWeaknesses:   detectedWeaknesses,
		ExecutedInstructions: executedInstructions,
		DeltaCoverage:        gofakeit.Uint64(),
		DeltaMinDistance:     gofakeit.Uint64(),
		Status:               common.RandomChoice([]common.TransactionStatus{common.TRANSACTION_DONE, common.TRANSACTION_RUNNING, common.TRANSACTION_CREATED}),
	}
}
