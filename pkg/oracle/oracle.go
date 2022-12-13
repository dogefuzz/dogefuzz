package oracle

import (
	"errors"

	"github.com/dogefuzz/dogefuzz/pkg/common"
)

type Oracle interface {
	Name() common.OracleType
	Detect(snapshot EventsSnapshot) bool
}

var ErrOracleDoesntExist = errors.New("oracle doesn't exist")

func GetOracles(oracleNames []common.OracleType) []Oracle {
	oracles := make([]Oracle, len(oracleNames))

	for _, oracleName := range oracleNames {
		oracle, err := GetOracleFromName(oracleName)
		if err != nil {
			continue
		}
		oracles = append(oracles, oracle)
	}
	return oracles
}

func GetOracleFromName(name common.OracleType) (Oracle, error) {
	switch name {
	case common.DELEGATE_ORACLE:
		return DelegateOracle{}, nil
	case common.EXCEPTION_DISORDER_ORACLE:
		return ExceptionDisorderOracle{}, nil
	case common.GASLESS_SEND_ORACLE:
		return GaslessSendOracle{}, nil
	case common.NUMBER_DEPENDENCY_ORACLE:
		return NumberDependencyOracle{}, nil
	case common.REENTRANCY_ORACLE:
		return ReentrancyOracle{}, nil
	case common.TIMESTAMP_DEPENDENCY_ORACLE:
		return TimestampDependencyOracle{}, nil
	}
	return nil, ErrOracleDoesntExist
}
