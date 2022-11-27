package common

func ConvertStringArrayToInterfaceArray(arr []string) []interface{} {
	output := make([]interface{}, len(arr))
	for value := range arr {
		output = append(output, value)
	}
	return output
}
