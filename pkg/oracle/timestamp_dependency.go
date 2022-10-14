package oracle

type TimestampDependencyOracle struct{}

func (o TimestampDependencyOracle) Name() string {
	return TIMESTAMP_DEPENDENCY_ORACLE
}

func (o TimestampDependencyOracle) Detect(snapshot EventsSnapshot) bool {
	return snapshot.Timestamp && (snapshot.StorageChanged || snapshot.EtherTransfer || snapshot.SendOp)
}
