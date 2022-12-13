package oracle

import "github.com/dogefuzz/dogefuzz/pkg/common"

type NumberDependencyOracle struct{}

func (o NumberDependencyOracle) Name() common.OracleType {
	return common.NUMBER_DEPENDENCY_ORACLE
}

func (o NumberDependencyOracle) Detect(snapshot EventsSnapshot) bool {
	return snapshot.BlockNumber && (snapshot.StorageChanged || snapshot.EtherTransfer || snapshot.SendOp)
}
