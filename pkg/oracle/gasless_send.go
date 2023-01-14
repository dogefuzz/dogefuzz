package oracle

import "github.com/dogefuzz/dogefuzz/pkg/common"

type GaslessSendOracle struct{}

func (o GaslessSendOracle) Name() common.OracleType {
	return common.GASLESS_SEND_ORACLE
}

func (o GaslessSendOracle) Detect(snapshot common.EventsSnapshot) bool {
	return snapshot.GaslessSend
}
