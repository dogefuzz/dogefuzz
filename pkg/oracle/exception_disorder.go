package oracle

import "github.com/dogefuzz/dogefuzz/pkg/common"

type ExceptionDisorderOracle struct{}

func (o ExceptionDisorderOracle) Name() common.OracleType {
	return common.EXCEPTION_DISORDER_ORACLE
}

func (o ExceptionDisorderOracle) Detect(snapshot common.EventsSnapshot) bool {
	return snapshot.ExceptionDisorder
}
