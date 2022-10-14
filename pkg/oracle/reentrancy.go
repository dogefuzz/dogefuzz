package oracle

type ReentrancyOracle struct{}

func (o ReentrancyOracle) Name() string {
	return REENTRANCY_ORACLE
}

func (o ReentrancyOracle) Detect(snapshot EventsSnapshot) bool {
	return snapshot.Reentrancy && (snapshot.StorageChanged || snapshot.EtherTransfer || snapshot.SendOp)
}
