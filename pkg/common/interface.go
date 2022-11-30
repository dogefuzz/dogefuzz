package common

func ConvertStringArrayToInterfaceArray(arr []string) []interface{} {
	output := make([]interface{}, len(arr))
	for idx := 0; idx < len(arr); idx++ {
		output[idx] = arr[idx]
	}
	return output
}
