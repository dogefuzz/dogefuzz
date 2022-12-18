package fuzz

import (
	"errors"

	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

var (
	ErrFuzzerTypeNotFound       = errors.New("fuzzer type not found")
	ErrFuzzerTypeNotImplemented = errors.New("fuzzer type not implemented")
)

type Fuzzer interface {
	GenerateInput(method abi.Method) []interface{}
}

func CreateFuzzer(typ common.FuzzingType) (Fuzzer, error) {
	switch typ {
	case common.BLACKBOX_FUZZING:
		return newBlackboxFuzzer(), nil
	case common.GREYBOX_FUZZING:
		return nil, ErrFuzzerTypeNotImplemented
	case common.DIRECTED_GREYBOX_FUZZING:
		return nil, ErrFuzzerTypeNotImplemented
	default:
		return nil, ErrFuzzerTypeNotFound
	}
}
