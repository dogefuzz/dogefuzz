package interfaces

import (
	"github.com/dogefuzz/dogefuzz/pkg/common"
)

type Oracle interface {
	Name() common.OracleType
	Detect(snapshot common.EventsSnapshot) bool
}
