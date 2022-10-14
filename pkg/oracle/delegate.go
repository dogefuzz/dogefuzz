package oracle

type DelegateOracle struct{}

func (o DelegateOracle) Name() string {
	return DELEGATE_ORACLE
}

func (o DelegateOracle) Detect(snapshot EventsSnapshot) bool {
	return snapshot.Delegate
}
