package oracle

import "github.com/gongbell/contractfuzzer/server/model"

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

func NewEventsSnapshot(oracles []model.OracleEvent) EventsSnapshot {
	snapshot := EventsSnapshot{}
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
