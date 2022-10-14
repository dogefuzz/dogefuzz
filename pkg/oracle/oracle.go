package oracle

import "errors"

type Oracle interface {
	Name() string
	Detect(snapshot EventsSnapshot) bool
}

var ErrOracleDoesntExist = errors.New("oracle doesn't exist")

func GetOracles(oracleNames []string) []Oracle {
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

func GetOracleFromName(name string) (Oracle, error) {
	switch name {
	case DELEGATE:
		return DelegateOracle{}, nil
	case EXCEPTION_DISORDER_ORACLE:
		return ExceptionDisorderOracle{}, nil
	case GASLESS_SEND_ORACLE:
		return GaslessSendOracle{}, nil
	case NUMBER_DEPENDENCY_ORACLE:
		return NumberDependencyOracle{}, nil
	case REENTRANCY_ORACLE:
		return ReentrancyOracle{}, nil
	case TIMESTAMP_DEPENDENCY_ORACLE:
		return TimestampDependencyOracle{}, nil
	}
	return nil, ErrOracleDoesntExist
}
