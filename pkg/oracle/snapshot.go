package oracle

import (
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
		if oracle == CALLFAILED {
			snapshot.CallFailed = true
		}
		if oracle == REENTRANCY {
			snapshot.Reentrancy = true
		}
		if oracle == REPEATED {
			snapshot.RepeatedCall = true
		}
		if oracle == ETHERTRANSFER {
			snapshot.EtherTransfer = true
		}
		if oracle == ETHERTRANSFERFAILED {
			snapshot.EtherTransferFailed = true
		}
		if oracle == CALLETHERETRANSFERFAILED {
			snapshot.CallEtherFailed = true
		}
		if oracle == GASLESSSEND {
			snapshot.GaslessSend = true
		}
		if oracle == DELEGATE {
			snapshot.Delegate = true
		}
		if oracle == EXCEPTIONDISORDER {
			snapshot.ExceptionDisorder = true
		}
		if oracle == SENDOP {
			snapshot.SendOp = true
		}
		if oracle == CALLOP {
			snapshot.CallOp = true
		}
		if oracle == CALLEXCEPTION {
			snapshot.CallException = true
		}
		if oracle == UNKNOWCALL {
			snapshot.UnknowCall = true
		}
		if oracle == STORAGECHANGE {
			snapshot.StorageChanged = true
		}
		if oracle == TIMESTAMP {
			snapshot.Timestamp = true
		}
		if oracle == BLOCKHAHSH {
			snapshot.BlockHash = true
		}
		if oracle == BLOCKNUMBER {
			snapshot.BlockNumber = true
		}
	}
	return snapshot
}
