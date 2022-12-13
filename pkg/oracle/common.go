package oracle

func GetFunctionNameFromInput(input string) string {
	if len(input) >= 0 {
		return input[:8]
	} else {
		return input
	}
}
