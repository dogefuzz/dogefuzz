package listener

import (
	"github.com/dogefuzz/dogefuzz/config"
	"github.com/dogefuzz/dogefuzz/pkg/bus"
	"github.com/dogefuzz/dogefuzz/pkg/interfaces"
	"go.uber.org/zap"
)

type Env interface {
	Logger() *zap.Logger
	Client() interfaces.HttpClient
	Config() *config.Config

	TaskStartTopic() interfaces.Topic[bus.TaskStartEvent]
	TaskFinishTopic() interfaces.Topic[bus.TaskFinishEvent]
	TaskInputRequestTopic() interfaces.Topic[bus.TaskInputRequestEvent]
	InstrumentExecutionTopic() interfaces.Topic[bus.InstrumentExecutionEvent]

	TaskService() interfaces.TaskService
	GethService() interfaces.GethService
	VandalService() interfaces.VandalService
	ContractService() interfaces.ContractService
	FunctionService() interfaces.FunctionService
	TransactionService() interfaces.TransactionService
	ReporterService() interfaces.ReporterService
	SolidityService() interfaces.SolidityService

	ContractMapper() interfaces.ContractMapper

	ContractDeployerListener() interfaces.Listener
	ExecutionAnalyticsListener() interfaces.Listener
	FuzzerListener() interfaces.Listener
	ReporterListener() interfaces.Listener

	Deployer() interfaces.Deployer
	Agent() interfaces.Agent

	FuzzerLeader() interfaces.FuzzerLeader
	BlackboxFuzzer() interfaces.Fuzzer
	GreyboxFuzzer() interfaces.Fuzzer
	DirectedGreyboxFuzzer() interfaces.Fuzzer
	PowerSchedule() interfaces.PowerSchedule
}
