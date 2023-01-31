package common

import (
	"math/big"
	"reflect"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

var (
	Uint8T      = reflect.TypeOf(uint8(0))
	Uint16T     = reflect.TypeOf(uint16(0))
	Uint32T     = reflect.TypeOf(uint32(0))
	Uint64T     = reflect.TypeOf(uint64(0))
	Int8T       = reflect.TypeOf(int8(0))
	Int16T      = reflect.TypeOf(int16(0))
	Int32T      = reflect.TypeOf(int32(0))
	Int64T      = reflect.TypeOf(int64(0))
	BigIntT     = reflect.TypeOf(&big.Int{})
	BoolT       = reflect.TypeOf(false)
	StringT     = reflect.TypeOf("")
	AddressT    = reflect.TypeOf(common.Address{})
	SliceT      = func(typ reflect.Type) reflect.Type { return reflect.SliceOf(typ) }
	ArrayT      = func(size int, typ reflect.Type) reflect.Type { return reflect.ArrayOf(size, typ) }
	FixedBytesT = func(size int) reflect.Type { return reflect.ArrayOf(size, reflect.TypeOf(byte(0))) }
	BytesT      = reflect.SliceOf(reflect.TypeOf(byte(0)))
)

type FuzzingType string

const (
	BLACKBOX_FUZZING         FuzzingType = "blackbox"
	GREYBOX_FUZZING          FuzzingType = "greybox"
	DIRECTED_GREYBOX_FUZZING FuzzingType = "directed_greybox"
)

type TaskStatus string

const (
	TASK_RUNNING TaskStatus = "running"
	TASK_DONE    TaskStatus = "done"
)

type TransactionStatus string

const (
	TRANSACTION_CREATED    TransactionStatus = "created"
	TRANSACTION_SEND_ERROR TransactionStatus = "send_error"
	TRANSACTION_RUNNING    TransactionStatus = "running"
	TRANSACTION_DONE       TransactionStatus = "done"
)

type OracleType string

const (
	DELEGATE_ORACLE             OracleType = "delegate"
	EXCEPTION_DISORDER_ORACLE   OracleType = "exception-disorder"
	GASLESS_SEND_ORACLE         OracleType = "gasless-send"
	NUMBER_DEPENDENCY_ORACLE    OracleType = "number-dependency"
	REENTRANCY_ORACLE           OracleType = "reentrancy"
	TIMESTAMP_DEPENDENCY_ORACLE OracleType = "timestamp-dependency"
)

type DistanceMap map[string]map[string]int64 // blockPC => instruction => distance

type PowerScheduleStrategy string

const (
	DISTANCE_BASED_STRATEGY PowerScheduleStrategy = "distance_based"
	COVERAGE_BASED_STRATEGY PowerScheduleStrategy = "coverage_based"
)

type TaskReport struct {
	TimeElapsed        time.Duration       `json:"timeElapsed"`
	ContractName       string              `json:"contractName"`
	Coverage           int64               `json:"coverage"`
	CoverageByTime     TimeSeriesData      `json:"coverageByTime"`
	MinDistance        int64               `json:"minDistance"`
	MinDistanceByTime  TimeSeriesData      `json:"minDistanceByTime"`
	Transactions       []TransactionReport `json:"transactions"`
	DetectedWeaknesses []string            `json:"detectedWeaknesses"`
}

type TimeSeriesData struct {
	X []time.Time `json:"x"`
	Y []int64     `json:"y"`
}

type TransactionReport struct {
	Timestamp            time.Time `json:"timestamp"`
	BlockchainHash       string    `json:"blockchainHash"`
	Inputs               []string  `json:"inputs"`
	DetectedWeaknesses   []string  `json:"detectedWeaknesses"`
	ExecutedInstructions []string  `json:"executedInstructions"`
	DeltaCoverage        int64     `json:"deltaCoverage"`
	DeltaMinDistance     int64     `json:"deltaMinDistance"`
}

type EventsSnapshot struct {
	CallFailed          bool
	Reentrancy          bool
	RepeatedCall        bool
	EtherTransfer       bool
	EtherTransferFailed bool
	CallEtherFailed     bool
	GaslessSend         bool
	Delegate            bool
	ExceptionDisorder   bool
	SendOp              bool
	CallOp              bool
	CallException       bool
	UnknowCall          bool
	StorageChanged      bool
	Timestamp           bool
	BlockHash           bool
	BlockNumber         bool
}

type TypeIdentifier string

type Block struct {
	PC               string                 `json:"pc"`
	Range            BlockRange             `json:"range"`
	Predecessors     []string               `json:"predecessors"`
	Successors       []string               `json:"successors"`
	EntryStack       []string               `json:"entryStack"`
	StackPops        uint64                 `json:"stackPops"`
	StackAdditions   []string               `json:"stackAdditions"`
	ExitStack        []string               `json:"exitStack"`
	Instructions     map[string]Instruction `json:"instructions"`
	InstructionOrder []string               `json:"instructionOrder"`
}

type BlockRange struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type Instruction struct {
	Op      string   `json:"op"`
	Args    []string `json:"args"`
	StackOp string   `json:"stackOp"`
}

type Function struct {
	Signature  string   `json:"signature"`
	EntryBlock string   `json:"entryBlock"`
	ExitBlock  string   `json:"exitBlock"`
	Body       []string `json:"body"`
}

type Seeds = map[TypeIdentifier][]string
