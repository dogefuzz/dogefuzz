package oracle

import (
	"strings"

	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/dogefuzz/dogefuzz/pkg/dto"
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

func NewEventsSnapshot(oracles []dto.OracleEvent) common.EventsSnapshot {
	snapshot := common.EventsSnapshot{}
	for _, oracle := range oracles {
		if strings.Contains(string(oracle), CALLFAILED) {
			snapshot.CallFailed = true
		}
		if strings.Contains(string(oracle), REENTRANCY) {
			snapshot.Reentrancy = true
		}
		if strings.Contains(string(oracle), REPEATED) {
			snapshot.RepeatedCall = true
		}
		if strings.Contains(string(oracle), ETHERTRANSFER) {
			snapshot.EtherTransfer = true
		}
		if strings.Contains(string(oracle), ETHERTRANSFERFAILED) {
			snapshot.EtherTransferFailed = true
		}
		if strings.Contains(string(oracle), CALLETHERETRANSFERFAILED) {
			snapshot.CallEtherFailed = true
		}
		if strings.Contains(string(oracle), GASLESSSEND) {
			snapshot.GaslessSend = true
		}
		if strings.Contains(string(oracle), DELEGATE) {
			snapshot.Delegate = true
		}
		if strings.Contains(string(oracle), EXCEPTIONDISORDER) {
			snapshot.ExceptionDisorder = true
		}
		if strings.Contains(string(oracle), SENDOP) {
			snapshot.SendOp = true
		}
		if strings.Contains(string(oracle), CALLOP) {
			snapshot.CallOp = true
		}
		if strings.Contains(string(oracle), CALLEXCEPTION) {
			snapshot.CallException = true
		}
		if strings.Contains(string(oracle), UNKNOWCALL) {
			snapshot.UnknowCall = true
		}
		if strings.Contains(string(oracle), STORAGECHANGE) {
			snapshot.StorageChanged = true
		}
		if strings.Contains(string(oracle), TIMESTAMP) {
			snapshot.Timestamp = true
		}
		if strings.Contains(string(oracle), BLOCKHAHSH) {
			snapshot.BlockHash = true
		}
		if strings.Contains(string(oracle), BLOCKNUMBER) {
			snapshot.BlockNumber = true
		}
	}
	return snapshot
}
