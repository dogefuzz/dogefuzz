package oracle

import (
	"fmt"

	"github.com/dogefuzz/dogefuzz/pkg/common"
)

type NumberDependencyOracle struct{}

func (o NumberDependencyOracle) Name() common.OracleType {
	return common.NUMBER_DEPENDENCY_ORACLE
}

func (o NumberDependencyOracle) Detect(snapshot common.EventsSnapshot) bool {
	fmt.Printf("%v && (%v || %v || %v)\n", snapshot.BlockNumber, snapshot.StorageChanged, snapshot.EtherTransfer, snapshot.SendOp)
	return snapshot.BlockNumber && (snapshot.StorageChanged || snapshot.EtherTransfer || snapshot.SendOp)
}
