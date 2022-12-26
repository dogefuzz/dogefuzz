package fuzz

import (
	"errors"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

var (
	ErrFuzzerTypeNotFound       = errors.New("fuzzer type not found")
	ErrFuzzerTypeNotImplemented = errors.New("fuzzer type not implemented")
)

type Fuzzer interface {
	GenerateInput(method abi.Method) []interface{}
}
