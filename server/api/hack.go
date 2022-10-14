package api

import (
	"encoding/csv"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gongbell/contractfuzzer/fuzz"
	"go.uber.org/zap"
)

const MAX_TIME uint64 = 10 * 60 // Ten minutes
const MAX_DURATION uint64 = 5

const (
	Danger_reentrancy          = "DR"
	Danger_exception_disorder  = "DE"
	Danger_delegate            = "DD"
	Danger_gasless_send        = "DGS"
	Danger_timestampdependency = "DT"
	Danger_numberdependency    = "DN"
	// Danger_freezingether = "DF"
)

const (
	CALLFAILED               = "HackerRootCallFailed"
	REENTRANCY               = "HackerReentrancy"
	REPEATED                 = "HackerRepeatedCall"
	ETHERTRANSFER            = "HackerEtherTransfer"
	ETHERTRANSFERFAILED      = "HackerEtherTransferFailed"
	CALLETHERETRANSFERFAILED = "HackerCallEtherTransferFailed"
	GASLESSSEND              = "HackerGaslessSend"
	DELEGATE                 = "HackerDelegateCallInfo"
	EXCEPTIONDISORDER        = "HackerExceptionDisorder"
	SENDOP                   = "HackerSendOpInfo"
	CALLOP                   = "HackerCallOpInfo"
	CALLEXCEPTION            = "HackerCallException"
	UNKNOWCALL               = "HackerUnknownCall"
	STORAGECHANGE            = "HackerStorageChanged"
	TIMESTAMP                = "HackerTimestampOp"
	BLOCKHAHSH               = "HackerBlockHashOp"
	BLOCKNUMBER              = "HackerNumberOp"
	FREEZINGETHER            = "HackerFreezingEther"
)

type HackAPI interface {
	Post(c *gin.Context)
}

type DefaultHackAPI struct {
	Logger                   *zap.Logger
	CountMutex               *sync.Mutex
	Mutex                    *sync.Mutex
	DecMutex                 *sync.Mutex
	TimeMutex                *sync.Mutex
	ReceiveCount             int
	AddrMapPath              string
	AddrMap                  map[string]string
	LogWriter                *os.File
	CountWriter              *os.File
	ContractOutputWriter     *os.File
	ReceiveCountWriter       *os.File
	ReentrancyWriter         *os.File
	ExceptionDisorderWriter  *os.File
	DelegateWriter           *os.File
	GaslessWriter            *os.File
	TimeDependencyWriter     *os.File
	NumberDependencyWriter   *os.File
	FreezingEtherWriter      *os.File
	StartTime                uint64
	TimeBound                uint64
	Result                   map[string]string
	Contract                 string
	LastContract             string
	Profile                  string
	LastProfile              string
	UnchangedDuration        uint64
	GroupCount               int
	CallFailedCount          int
	StorageChangedCount      int
	CallOpCount              int
	CallExceptionCount       int
	ExceptionDisorderCount   int
	EtherTransferCount       int
	EtherTransferFailedCount int
	DelegateCount            int
	GaslessSendCount         int
	ReentrancyCount          int
	CallEtherFailedCount     int
	RepeatedCallCount        int
	TimestampCount           int
	BlockHashCount           int
	BlockNumberCount         int
	SendOpCount              int
	UnknowCallCount          int
	FreezingEtherCount       int
	IsCallFailed             bool
	IsStorageChanged         bool
	IsCallOp                 bool
	IsCallException          bool
	IsExceptionDisorder      bool
	IsEtherTransfer          bool
	IsEtherTransferFailed    bool
	IsDelegate               bool
	IsGaslessSend            bool
	IsReentrancy             bool
	IsCallEtherFailed        bool
	IsRepeatedCall           bool
	IsTimestamp              bool
	IsBlockHash              bool
	IsBlockNumber            bool
	IsSendOp                 bool
	IsUnknowCall             bool
	// IsFreezingEther bool
}

func (api DefaultHackAPI) Init(logger *zap.Logger, addrMapPath string, reporterPath string) DefaultHackAPI {
	defer api.handleInitPanic()
	api.Logger = logger

	api.initFilesWithRWPermissions(reporterPath)
	if _, err := os.Stat(addrMapPath); errors.Is(err, os.ErrNotExist) {
		errorMsg := fmt.Sprintf("The provided address map file does not exist: %s", err)
		api.Logger.Panic(errorMsg)
		panic(errorMsg)
	}
	api.AddrMapPath = addrMapPath
	api.importAddrContractMap()

	api.Result = make(map[string]string)
	api.Contract = ""
	api.LastContract = ""
	api.Profile = ""
	api.LastProfile = ""
	api.UnchangedDuration = MAX_TIME
	api.ReceiveCount = 0
	api.CountMutex = new(sync.Mutex)
	api.Mutex = new(sync.Mutex)
	api.DecMutex = new(sync.Mutex)
	api.TimeMutex = new(sync.Mutex)
	api.GroupCount = 0

	api.CallFailedCount = 0
	api.StorageChangedCount = 0
	api.CallOpCount = 0
	api.CallExceptionCount = 0
	api.ExceptionDisorderCount = 0
	api.EtherTransferCount = 0
	api.EtherTransferFailedCount = 0
	api.DelegateCount = 0
	api.GaslessSendCount = 0
	api.ReentrancyCount = 0
	api.CallEtherFailedCount = 0
	api.RepeatedCallCount = 0
	api.TimestampCount = 0
	api.BlockHashCount = 0
	api.BlockNumberCount = 0
	api.SendOpCount = 0
	api.UnknowCallCount = 0
	api.FreezingEtherCount = 0
	api.IsCallFailed = false
	api.IsStorageChanged = false
	api.IsCallOp = false
	api.IsCallException = false
	api.IsExceptionDisorder = false
	api.IsEtherTransfer = false
	api.IsEtherTransferFailed = false
	api.IsDelegate = false
	api.IsGaslessSend = false
	api.IsReentrancy = false
	api.IsCallEtherFailed = false
	api.IsRepeatedCall = false
	api.IsTimestamp = false
	api.IsBlockHash = false
	api.IsBlockNumber = false
	api.IsSendOp = false
	api.IsUnknowCall = false
	// api.IsFreezingEther = false

	api.StartTime = uint64(time.Now().Unix())
	api.TimeBound = api.StartTime + MAX_TIME

	go api.startFuzzerTimer()
	go api.checkNewContractAndRestartFuzzerTimer()

	return api
}

func (api DefaultHackAPI) Post(c *gin.Context) {
	defer api.handlePanic(c)

	api.CountMutex.Lock()
	api.ReceiveCount++
	api.CountMutex.Unlock()

	api.ReceiveCountWriter.Write([]byte(strconv.Itoa(api.ReceiveCount)))

	oracles := c.Request.URL.Query().Get("oracles")
	profile := c.Request.URL.Query().Get("profile")
	txHash := c.Request.URL.Query().Get("txHash")
	api.Logger.Info(fmt.Sprintf("TX HASH: %s", txHash))
	api.countOracles(oracles)
	countList := fmt.Sprintf("%4d %4d %4d %4d %4d %4d %4d %4d %4d %4d %4d %4d %4d %4d %4d %4d %4d %4d\n ", api.CallFailedCount, api.StorageChangedCount, api.CallOpCount, api.CallExceptionCount, api.ExceptionDisorderCount, api.EtherTransferCount, api.EtherTransferFailedCount, api.DelegateCount, api.GaslessSendCount, api.ReentrancyCount, api.CallEtherFailedCount, api.RepeatedCallCount, api.TimestampCount, api.BlockHashCount, api.BlockNumberCount, api.SendOpCount, api.UnknowCallCount, api.FreezingEtherCount)
	api.CountWriter.Write([]byte(countList))
	msg := fmt.Sprintf("\n%s  %s\n", oracles, profile)
	api.LogWriter.Write([]byte(msg))
	api.Mutex.Lock()
	api.hackCount(oracles, profile)
	api.Mutex.Unlock()

	api.DecMutex.Lock()
	api.GroupCount--
	api.Logger.Info(fmt.Sprintf("Group count: %d", api.GroupCount))
	if api.GroupCount <= 0 {
		if compareProfiles(api.LastProfile, api.Profile) {
			api.UnchangedDuration--
			if api.UnchangedDuration == 0 {
				fuzz.G_stop <- true
			} else {
				fuzz.G_stop <- false
			}
		} else {
			fuzz.G_stop <- false
			api.UnchangedDuration = MAX_DURATION
		}
		api.LastProfile = api.Profile
	}
	api.DecMutex.Unlock()

	var resp struct {
		Message string `json:"message"`
	}
	resp.Message = fmt.Sprintf("GET params were:%s", c.Request.URL.Query())
	c.JSON(http.StatusOK, resp)
}

func (api DefaultHackAPI) initFilesWithRWPermissions(reporterPath string) {
	logFile := filepath.Join(reporterPath, "log.txt")
	api.LogWriter = api.openFileOrPanic(logFile)

	countFile := filepath.Join(reporterPath, "count.txt")
	api.CountWriter = api.openFileOrPanic(countFile)

	contractOutputFile := filepath.Join(reporterPath, "contract_fun_vulnerabilities.txt")
	api.ContractOutputWriter = api.openFileOrPanic(contractOutputFile)

	receiveCountFile := filepath.Join(reporterPath, "receive_count.txt")
	api.ReceiveCountWriter = api.openFileOrPanic(receiveCountFile)

	reentrancyFile := filepath.Join(reporterPath, "/bug/reentrancy_danger.list")
	api.ReentrancyWriter = api.openFileOrPanic(reentrancyFile)

	exceptionDisorderFile := filepath.Join(reporterPath, "/bug/exception_disorder.list")
	api.ExceptionDisorderWriter = api.openFileOrPanic(exceptionDisorderFile)

	delegateFile := filepath.Join(reporterPath, "/bug/delegate_danger.list")
	api.DelegateWriter = api.openFileOrPanic(delegateFile)

	gaslessFile := filepath.Join(reporterPath, "/bug/gasless_send.list")
	api.GaslessWriter = api.openFileOrPanic(gaslessFile)

	timeDependencyFile := filepath.Join(reporterPath, "/bug/time_dependency.list")
	api.TimeDependencyWriter = api.openFileOrPanic(timeDependencyFile)

	numberDependencyFile := filepath.Join(reporterPath, "/bug/number_dependency.list")
	api.NumberDependencyWriter = api.openFileOrPanic(numberDependencyFile)

	freezingEtherFile := filepath.Join(reporterPath, "/bug/freezing_ether.list")
	api.FreezingEtherWriter = api.openFileOrPanic(freezingEtherFile)
}

func (api DefaultHackAPI) importAddrContractMap() {
	f, err := os.Open(api.AddrMapPath)
	if err != nil {
		errorMsg := fmt.Sprintf("The address map file could not be opened: %s", err)
		api.Logger.Panic(errorMsg)
		panic(errorMsg)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		errorMsg := fmt.Sprintf("Error while reading the address map as a csv file: %s", err)
		api.Logger.Panic(errorMsg)
		panic(errorMsg)
	}

	api.AddrMap = make(map[string]string, 0)
	for _, line := range records {
		addr := line[0]
		contractName := line[1]
		api.AddrMap[strings.ToLower(addr)] = contractName
	}
}

func (api DefaultHackAPI) startFuzzerTimer() {
	for {
		api.Logger.Info(fmt.Sprintf("Current Time: %d -- Time Bound: %d", time.Now().Unix(), api.TimeBound))
		if uint64(time.Now().Unix()) > api.TimeBound {
			api.Logger.Info("TIMEOUT, stop.")
			fuzz.G_stop <- true
		}
		time.Sleep(60 * time.Second)
	}
}

func (api DefaultHackAPI) checkNewContractAndRestartFuzzerTimer() {
	for {
		<-fuzz.G_start
		api.TimeMutex.Lock()
		if api.Contract != api.LastContract {
			api.StartTime = uint64(time.Now().Unix())
			api.TimeBound = api.StartTime + MAX_TIME
			api.LastContract = api.Contract
			api.UnchangedDuration = MAX_DURATION
		}
		api.TimeMutex.Unlock()
		api.Logger.Info(fmt.Sprintf("Rand case scale: %d", fuzz.RAND_CASE_SCALE))
		api.GroupCount = 3 * fuzz.RAND_CASE_SCALE
		api.Logger.Info(fmt.Sprintf("Group count: %d", api.GroupCount))
		fuzz.G_sig_continue <- true
	}
}

func (api DefaultHackAPI) countOracles(oracles string) {
	if strings.Contains(oracles, CALLFAILED) {
		api.CallFailedCount++
	}
	if strings.Contains(oracles, REENTRANCY) {
		api.ReentrancyCount++
	}
	if strings.Contains(oracles, REPEATED) {
		api.RepeatedCallCount++
	}
	if strings.Contains(oracles, ETHERTRANSFER) {
		api.EtherTransferCount++
	}
	if strings.Contains(oracles, ETHERTRANSFERFAILED) {
		api.EtherTransferFailedCount++
	}
	if strings.Contains(oracles, CALLETHERETRANSFERFAILED) {
		api.CallEtherFailedCount++
	}
	if strings.Contains(oracles, GASLESSSEND) {
		api.GaslessSendCount++
	}
	if strings.Contains(oracles, DELEGATE) {
		api.DelegateCount++
	}
	if strings.Contains(oracles, EXCEPTIONDISORDER) {
		api.ExceptionDisorderCount++
	}
	if strings.Contains(oracles, SENDOP) {
		api.SendOpCount++
	}
	if strings.Contains(oracles, CALLOP) {
		api.CallOpCount++
	}
	if strings.Contains(oracles, CALLEXCEPTION) {
		api.CallExceptionCount++
	}
	if strings.Contains(oracles, UNKNOWCALL) {
		api.UnknowCallCount++
	}
	if strings.Contains(oracles, STORAGECHANGE) {
		api.StorageChangedCount++
	}
	if strings.Contains(oracles, TIMESTAMP) {
		api.TimestampCount++
	}
	if strings.Contains(oracles, BLOCKHAHSH) {
		api.BlockHashCount++
	}
	if strings.Contains(oracles, BLOCKNUMBER) {
		api.BlockNumberCount++
	}
	if strings.Contains(oracles, FREEZINGETHER) {
		api.FreezingEtherCount++
	}
}

func (api DefaultHackAPI) hackCount(oracles, profiles string) {
	caller := strings.Split(strings.Split(profiles, ",")[0], "caller:")[1]
	newContract := strings.Split(strings.Split(profiles, ",")[1], "callee:")[1]
	value := strings.Split(strings.Split(profiles, ",")[2], "value:")[1]
	input := strings.Split(strings.Split(strings.Split(profiles, ",")[4], ":")[1], "}")[0]

	if newContract != api.Contract {
		api.output()
		api.Contract = newContract
		api.Result = make(map[string]string)
		api.Profile = ""
	}
	api.resetVulnerabilities()
	if strings.Contains(oracles, CALLFAILED) {
		api.IsCallFailed = true
	}
	if strings.Contains(oracles, REENTRANCY) {
		api.IsReentrancy = true
	}
	if strings.Contains(oracles, REPEATED) {
		api.IsRepeatedCall = true
	}
	if strings.Contains(oracles, ETHERTRANSFER) {
		api.IsEtherTransfer = true
	}
	if strings.Contains(oracles, ETHERTRANSFERFAILED) {
		api.IsEtherTransferFailed = true
	}
	if strings.Contains(oracles, CALLETHERETRANSFERFAILED) {
		api.IsCallEtherFailed = true
	}
	if strings.Contains(oracles, GASLESSSEND) {
		api.IsGaslessSend = true
	}
	if strings.Contains(oracles, DELEGATE) {
		//isGaslessSend = true; ERROR：校正gaslessSend数据
		api.IsDelegate = true
	}

	if strings.Contains(oracles, EXCEPTIONDISORDER) {
		api.IsExceptionDisorder = true
	}
	if strings.Contains(oracles, SENDOP) {
		api.IsSendOp = true
	}
	if strings.Contains(oracles, CALLOP) {
		api.IsCallOp = true
	}
	if strings.Contains(oracles, CALLEXCEPTION) {
		api.IsCallException = true
	}
	if strings.Contains(oracles, UNKNOWCALL) {
		api.IsUnknowCall = true
	}
	if strings.Contains(oracles, STORAGECHANGE) {
		api.IsStorageChanged = true
	}
	if strings.Contains(oracles, TIMESTAMP) {
		api.IsTimestamp = true
	}
	if strings.Contains(oracles, BLOCKHAHSH) {
		api.IsBlockHash = true
	}
	if strings.Contains(oracles, BLOCKNUMBER) {
		api.IsBlockNumber = true
	}
	api.hackCountFunc(caller, api.Contract, value, input)
}

func (api DefaultHackAPI) hackCountFunc(caller, callee, value, input string) {
	fun := input
	if len(input) >= 8 {
		fun = input[:8]
	}
	if api.IsReentrancy && (api.IsStorageChanged || api.IsEtherTransfer || api.IsSendOp) {
		if dangers, found := api.Result[fun]; found {
			if !strings.Contains(dangers, Danger_reentrancy) {
				api.Result[fun] = api.Result[fun] + " " + Danger_reentrancy
			}
		} else {
			api.Result[fun] = Danger_reentrancy
		}
		if !strings.Contains(api.Profile, Danger_reentrancy) {
			api.Profile += " " + Danger_reentrancy
		}
	}
	if api.IsExceptionDisorder {
		if dangers, found := api.Result[fun]; found {
			if !strings.Contains(dangers, Danger_exception_disorder) {
				api.Result[fun] = api.Result[fun] + " " + Danger_exception_disorder
			}
		} else {
			api.Result[fun] = Danger_exception_disorder
		}
		if !strings.Contains(api.Profile, Danger_exception_disorder) {
			api.Profile += " " + Danger_exception_disorder
		}
	}
	if api.IsDelegate {
		if dangers, found := api.Result[fun]; found {
			if !strings.Contains(dangers, Danger_delegate) {
				api.Result[fun] = api.Result[fun] + " " + Danger_delegate
			}
		} else {
			api.Result[fun] = Danger_delegate
		}
		if !strings.Contains(api.Profile, Danger_delegate) {
			api.Profile += " " + Danger_delegate
		}
	}
	if api.IsGaslessSend {
		if dangers, found := api.Result[fun]; found {
			if !strings.Contains(dangers, Danger_gasless_send) {
				api.Result[fun] = api.Result[fun] + " " + Danger_gasless_send
			}
		} else {
			api.Result[fun] = Danger_gasless_send
		}
		if !strings.Contains(api.Profile, Danger_gasless_send) {
			api.Profile += " " + Danger_gasless_send
		}
	}
	if api.IsTimestamp && (api.IsStorageChanged || api.IsEtherTransfer || api.IsSendOp) {
		if dangers, found := api.Result[fun]; found {
			if !strings.Contains(dangers, Danger_timestampdependency) {
				api.Result[fun] = api.Result[fun] + " " + Danger_timestampdependency
			}
		} else {
			api.Result[fun] = Danger_timestampdependency
		}
		if !strings.Contains(api.Profile, Danger_timestampdependency) {
			api.Profile += " " + Danger_timestampdependency
		}
	}
	if api.IsBlockNumber && (api.IsStorageChanged || api.IsEtherTransfer || api.IsSendOp) {
		if dangers, found := api.Result[fun]; found {
			if !strings.Contains(dangers, Danger_numberdependency) {
				api.Result[fun] = api.Result[fun] + " " + Danger_numberdependency
			}
		} else {
			api.Result[fun] = Danger_numberdependency
		}
		if !strings.Contains(api.Profile, Danger_numberdependency) {
			api.Profile += " " + Danger_numberdependency
		}
	}
}

func (api DefaultHackAPI) output() {
	if api.Profile == "" {
		return
	}
	packetSplit := "\n***********************"
	api.ContractOutputWriter.Write([]byte(packetSplit))
	packetInfo := "\n" + api.AddrMap[strings.ToLower(api.Contract[1:len(api.Contract)-1])]
	packetInfo += "\n" + api.Profile
	for fun := range api.Result {
		packetInfo += "\n" + fun + ": " + api.Result[fun]
	}
	if len(api.Result) > 0 {
		api.ContractOutputWriter.Write([]byte(packetInfo))
		api.Logger.Info(fmt.Sprintf("PackageInfo: %s", packetInfo))
	}
	if strings.Contains(api.Profile, Danger_reentrancy) {
		api.ReentrancyWriter.Write([]byte(api.AddrMap[strings.ToLower(api.Contract[1:len(api.Contract)-1])]))
		api.ReentrancyWriter.Write([]byte("\n"))
	}
	if strings.Contains(api.Profile, Danger_exception_disorder) {
		api.ExceptionDisorderWriter.Write([]byte(api.AddrMap[strings.ToLower(api.Contract[1:len(api.Contract)-1])]))
		api.ExceptionDisorderWriter.Write([]byte("\n"))
	}
	if strings.Contains(api.Profile, Danger_delegate) {
		api.DelegateWriter.Write([]byte(api.AddrMap[strings.ToLower(api.Contract[1:len(api.Contract)-1])]))
		api.DelegateWriter.Write([]byte("\n"))
	}
	if strings.Contains(api.Profile, Danger_gasless_send) {
		api.GaslessWriter.Write([]byte(api.AddrMap[strings.ToLower(api.Contract[1:len(api.Contract)-1])]))
		api.GaslessWriter.Write([]byte("\n"))
	}
	if strings.Contains(api.Profile, Danger_timestampdependency) {
		api.TimeDependencyWriter.Write([]byte(api.AddrMap[strings.ToLower(api.Contract[1:len(api.Contract)-1])]))
		api.TimeDependencyWriter.Write([]byte("\n"))
	}
	if strings.Contains(api.Profile, Danger_numberdependency) {
		api.NumberDependencyWriter.Write([]byte(api.AddrMap[strings.ToLower(api.Contract[1:len(api.Contract)-1])]))
		api.NumberDependencyWriter.Write([]byte("\n"))
	}
	// if strings.Contains(api.Profile, Danger_freezingether) {
	// 	api.FreezingEtherWriter.Write([]byte(api.AddrMap[strings.ToLower(api.Contract[1:len(api.Contract)-1])]))
	// 	api.FreezingEtherWriter.Write([]byte("\n"))
	// }
}

func (api DefaultHackAPI) resetVulnerabilities() {
	api.IsCallFailed = false
	api.IsStorageChanged = false
	api.IsCallOp = false
	api.IsCallException = false
	api.IsExceptionDisorder = false
	api.IsEtherTransfer = false
	api.IsEtherTransferFailed = false
	api.IsDelegate = false
	api.IsGaslessSend = false
	api.IsReentrancy = false
	api.IsCallEtherFailed = false
	api.IsRepeatedCall = false
	api.IsTimestamp = false
	api.IsBlockHash = false
	api.IsBlockNumber = false
	api.IsSendOp = false
	api.IsUnknowCall = false
	// api.IsFreezingEther = false
}

func (api DefaultHackAPI) openFileOrPanic(path string) *os.File {
	countWriter, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		errorMsg := fmt.Sprintf("The file could not be created: %s", err)
		api.Logger.Panic(errorMsg)
		panic(errorMsg)
	}
	return countWriter
}

func (api DefaultHackAPI) handlePanic(c *gin.Context) {
	if err := recover(); err != nil {
		api.Logger.Panic(err.(string))
		c.AbortWithError(500, errors.New(err.(string)))
	}
}

func (api DefaultHackAPI) handleInitPanic() {
	if err := recover(); err != nil {
		api.Logger.Info(err.(string))
		for i := 0; i < 10; i++ {
			funcName, file, line, ok := runtime.Caller(i)
			if ok {
				api.Logger.Info(fmt.Sprintf("frame %v:[func:%v,file:%v,line:%v]\n", i, runtime.FuncForPC(funcName).Name(), file, line))
			}
		}
	}
}

func compareProfiles(oldprofile, profile string) bool {
	return strings.EqualFold(oldprofile, profile)
}
