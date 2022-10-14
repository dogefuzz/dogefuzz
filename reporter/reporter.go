package reporter

import "github.com/dogefuzz/dogefuzz/env"

type Reporter interface {
}

type DefaultReporter struct {
	environment env.Environment
}

func (r DefaultReporter) Init(environment env.Environment) {
	r.environment = environment
}
