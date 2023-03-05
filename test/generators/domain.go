package generators

import (
	"encoding/json"
	"strconv"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/dogefuzz/dogefuzz/entities"
	"github.com/dogefuzz/dogefuzz/pkg/common"
)

func ContractGen() *entities.Contract {
	cfg, _ := json.Marshal(CFGGen())
	distanceMap, _ := json.Marshal(DistanceMapGen())
	return &entities.Contract{
		Id:                 gofakeit.LetterN(255),
		TaskId:             gofakeit.LetterN(255),
		Address:            SmartContractGen(),
		Source:             gofakeit.LetterN(255),
		DeploymentBytecode: gofakeit.LetterN(255),
		AbiDefinition:      gofakeit.LetterN(255),
		Name:               gofakeit.Name(),
		CFG:                string(cfg),
		DistanceMap:        string(distanceMap),
	}
}

func FunctionGen() *entities.Function {
	return &entities.Function{
		Name:          gofakeit.Name(),
		NumberOfArgs:  gofakeit.Int64(),
		IsConstructor: gofakeit.Bool(),
		ContractId:    gofakeit.LetterN(255),
	}
}

func TaskGen() *entities.Task {
	return &entities.Task{
		Arguments:                      gofakeit.LetterN(255),
		StartTime:                      gofakeit.Date(),
		Expiration:                     gofakeit.Date(),
		Detectors:                      gofakeit.LetterN(255),
		FuzzingType:                    common.RandomChoice([]common.FuzzingType{common.BLACKBOX_FUZZING, common.GREYBOX_FUZZING, common.DIRECTED_GREYBOX_FUZZING}),
		AggregatedExecutedInstructions: gofakeit.LetterN(255),
		Status:                         common.RandomChoice([]common.TaskStatus{common.TASK_RUNNING, common.TASK_DONE}),
	}
}

func TransactionGen() *entities.Transaction {
	return &entities.Transaction{
		Timestamp:            gofakeit.Date(),
		BlockchainHash:       gofakeit.LetterN(255),
		TaskId:               gofakeit.LetterN(255),
		FunctionId:           gofakeit.LetterN(255),
		Inputs:               gofakeit.LetterN(255),
		DetectedWeaknesses:   gofakeit.LetterN(255),
		ExecutedInstructions: gofakeit.LetterN(255),
		DeltaCoverage:        strconv.FormatUint(gofakeit.Uint64(), 10),
		DeltaMinDistance:     strconv.FormatUint(gofakeit.Uint64(), 10),
		Status:               common.RandomChoice([]common.TransactionStatus{common.TRANSACTION_CREATED, common.TRANSACTION_RUNNING, common.TRANSACTION_SEND_ERROR, common.TRANSACTION_DONE}),
	}
}
