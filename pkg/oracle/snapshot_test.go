package oracle

import (
	"testing"

	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/dogefuzz/dogefuzz/pkg/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SnapshotTestSuite struct {
	suite.Suite
}

func (suite *SnapshotTestSuite) TestEventsSnapshotCreationWithTrueValues() {
	oracles := []dto.OracleEvent{
		"HackerRootCallFailed",
		"HackerReentrancy",
		"HackerRepeatedCall",
		"HackerEtherTransfer",
		"HackerEtherTransferFailed",
		"HackerCallEtherTransferFailed",
		"HackerGaslessSend",
		"HackerDelegateCallInfo",
		"HackerExceptionDisorder",
		"HackerSendOpInfo",
		"HackerCallOpInfo",
		"HackerCallException",
		"HackerUnknownCall",
		"HackerStorageChanged",
		"HackerTimestampOp",
		"HackerBlockHashOp",
		"HackerNumberOp",
		"HackerFreezingEther",
	}

	expectedSnapshot := common.EventsSnapshot{
		CallFailed:          true,
		Reentrancy:          true,
		RepeatedCall:        true,
		EtherTransfer:       true,
		EtherTransferFailed: true,
		CallEtherFailed:     true,
		GaslessSend:         true,
		Delegate:            true,
		ExceptionDisorder:   true,
		SendOp:              true,
		CallOp:              true,
		CallException:       true,
		UnknowCall:          true,
		StorageChanged:      true,
		Timestamp:           true,
		BlockHash:           true,
		BlockNumber:         true,
	}

	assert.Equal(suite.T(), expectedSnapshot, NewEventsSnapshot(oracles), "delegate call didn't detect weakness")
}

func (suite *SnapshotTestSuite) TestEventsSnapshotCreationWithFalseValues() {
	oracles := []dto.OracleEvent{}

	expectedSnapshot := common.EventsSnapshot{
		CallFailed:          false,
		Reentrancy:          false,
		RepeatedCall:        false,
		EtherTransfer:       false,
		EtherTransferFailed: false,
		CallEtherFailed:     false,
		GaslessSend:         false,
		Delegate:            false,
		ExceptionDisorder:   false,
		SendOp:              false,
		CallOp:              false,
		CallException:       false,
		UnknowCall:          false,
		StorageChanged:      false,
		Timestamp:           false,
		BlockHash:           false,
		BlockNumber:         false,
	}

	assert.Equal(suite.T(), expectedSnapshot, NewEventsSnapshot(oracles), "delegate call didn't detect weakness")
}

func TestSnapshotTestSuite(t *testing.T) {
	suite.Run(t, new(SnapshotTestSuite))
}
