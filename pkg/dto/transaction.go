package dto

import (
	"time"

	"github.com/dogefuzz/dogefuzz/pkg/common"
)

type NewTransactionDTO struct {
	Timestamp  time.Time                `json:"timestamp"`
	TaskId     string                   `json:"taskId"`
	FunctionId string                   `json:"functionId"`
	Inputs     []string                 `json:"inputs"`
	Status     common.TransactionStatus `json:"status"`
}

type TransactionDTO struct {
	Id                   string                   `json:"id"`
	Timestamp            time.Time                `json:"timestamp"`
	BlockchainHash       string                   `json:"blockchainHash"`
	TaskId               string                   `json:"taskId"`
	FunctionId           string                   `json:"functionId"`
	Inputs               []string                 `json:"inputs"`
	DetectedWeaknesses   []string                 `json:"detectedWeaknesses"`
	ExecutedInstructions []string                 `json:"executedInstructions"`
	DeltaCoverage        int64                    `json:"deltaCoverage"`
	DeltaMinDistance     int64                    `json:"deltaMinDistance"`
	Status               common.TransactionStatus `json:"status"`
}

type NewExecutionDTO struct {
	Name         string   `json:"name"`
	Input        string   `json:"input"`
	Instructions []uint64 `json:"instructions"`
	TxHash       string   `json:"txHash"`
}

type NewWeaknessDTO struct {
	OracleEvents []OracleEvent `json:"oracleEvents"`
	Execution    Execution     `json:"execution"`
	TxHash       string        `json:"txHash"`
}

type OracleEvent string

// type Profile string

type Execution struct {
	Metadata    ExecutionMetadata `json:"metadata"`
	CallsLength int               `json:"callsLength"`
	Trace       ExecutionTrace    `json:"trace"`
}

type ExecutionMetadata struct {
	Caller string `json:"caller"`
	Callee string `json:"callee"`
	Value  string `json:"value"`
	Gas    string `json:"gas"`
	Input  string `json:"input"`
}

type ExecutionTrace []string
