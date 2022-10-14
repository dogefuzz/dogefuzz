package event

import "github.com/gongbell/contractfuzzer/db/domain"

type CoverageEvent struct {
	Input        string
	Instructions []uint64
	Transaction  domain.Transaction
}
