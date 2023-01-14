package oracle

import "github.com/dogefuzz/dogefuzz/pkg/common"

type ReentrancyOracle struct{}

func (o ReentrancyOracle) Name() common.OracleType {
	return common.REENTRANCY_ORACLE
}

func (o ReentrancyOracle) Detect(snapshot common.EventsSnapshot) bool {
	return snapshot.Reentrancy && (snapshot.StorageChanged || snapshot.EtherTransfer || snapshot.SendOp)
}
