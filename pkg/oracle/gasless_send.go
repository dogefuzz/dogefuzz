package oracle

type GaslessSendOracle struct{}

func (o GaslessSendOracle) Name() string {
	return GASLESS_SEND_ORACLE
}

func (o GaslessSendOracle) Detect(snapshot EventsSnapshot) bool {
	return snapshot.GaslessSend
}
