package fuzz

import (
	"github.com/dogefuzz/dogefuzz/config"
	"github.com/dogefuzz/dogefuzz/service"
)

type env interface {
	Config() *config.Config
	BlackboxFuzzer() Fuzzer
	GreyboxFuzzer() Fuzzer
	DirectedGreyboxFuzzer() Fuzzer
	PowerSchedule() PowerSchedule
	TransactionService() service.TransactionService
}
