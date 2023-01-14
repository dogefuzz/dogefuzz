package fuzz

import (
	"errors"
)

var (
	ErrFuzzerTypeNotFound       = errors.New("fuzzer type not found")
	ErrFuzzerTypeNotImplemented = errors.New("fuzzer type not implemented")
)
