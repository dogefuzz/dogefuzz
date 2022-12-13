package common

import (
	"math/big"
	"reflect"
)

var (
	Uint8T  = reflect.TypeOf(uint8(0))
	Uint16T = reflect.TypeOf(uint16(0))
	Uint32T = reflect.TypeOf(uint32(0))
	Uint64T = reflect.TypeOf(uint64(0))
	Int8T   = reflect.TypeOf(int8(0))
	Int16T  = reflect.TypeOf(int16(0))
	Int32T  = reflect.TypeOf(int32(0))
	Int64T  = reflect.TypeOf(int64(0))
	BigIntT = reflect.TypeOf(&big.Int{})
	BoolT   = reflect.TypeOf(false)
	StringT = reflect.TypeOf("")
	SliceT  = func(typ reflect.Type) reflect.Type { return reflect.SliceOf(typ) }
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
	TRANSACTION_RUNNING TransactionStatus = "running"
	TRANSACTION_DONE    TransactionStatus = "done"
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
