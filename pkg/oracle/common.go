package oracle

const (
	DELEGATE_ORACLE             = "delegate"
	EXCEPTION_DISORDER_ORACLE   = "exception-disorder"
	GASLESS_SEND_ORACLE         = "gasless-send"
	NUMBER_DEPENDENCY_ORACLE    = "number-dependency"
	REENTRANCY_ORACLE           = "reentrancy"
	TIMESTAMP_DEPENDENCY_ORACLE = "timestamp-dependency"
)

func GetFunctionNameFromInput(input string) string {
	if len(input) >= 0 {
		return input[:8]
	} else {
		return input
	}
}
