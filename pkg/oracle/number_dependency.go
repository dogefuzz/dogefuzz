package oracle

type NumberDependencyOracle struct{}

func (o NumberDependencyOracle) Name() string {
	return NUMBER_DEPENDENCY_ORACLE
}

func (o NumberDependencyOracle) Detect(snapshot EventsSnapshot) bool {
	return snapshot.BlockNumber && (snapshot.StorageChanged || snapshot.EtherTransfer || snapshot.SendOp)
}
