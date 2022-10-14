package oracle

type ExceptionDisorderOracle struct{}

func (o ExceptionDisorderOracle) Name() string {
	return EXCEPTION_DISORDER_ORACLE
}

func (o ExceptionDisorderOracle) Detect(snapshot EventsSnapshot) bool {
	return snapshot.ExceptionDisorder
}
