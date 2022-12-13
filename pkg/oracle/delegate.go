package oracle

import "github.com/dogefuzz/dogefuzz/pkg/common"

type DelegateOracle struct{}

func (o DelegateOracle) Name() common.OracleType {
	return common.DELEGATE_ORACLE
}

func (o DelegateOracle) Detect(snapshot EventsSnapshot) bool {
	return snapshot.Delegate
}
