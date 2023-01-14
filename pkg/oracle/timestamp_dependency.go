package oracle

import "github.com/dogefuzz/dogefuzz/pkg/common"

type TimestampDependencyOracle struct{}

func (o TimestampDependencyOracle) Name() common.OracleType {
	return common.TIMESTAMP_DEPENDENCY_ORACLE
}

func (o TimestampDependencyOracle) Detect(snapshot common.EventsSnapshot) bool {
	return snapshot.Timestamp && (snapshot.StorageChanged || snapshot.EtherTransfer || snapshot.SendOp)
}
