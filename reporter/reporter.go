package reporter

type Reporter interface {
}

type reporter struct {
}

func NewReporter() *reporter {
	return &reporter{}
}
