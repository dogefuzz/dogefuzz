package event

import "github.com/dogefuzz/dogefuzz/db/domain"

type CoverageEvent struct {
	Input        string
	Instructions []uint64
	Transaction  domain.Transaction
}
