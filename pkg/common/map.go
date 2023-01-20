package common

func GetFirstStringKeyFromMap[V any](input map[string]V) string {
	var firstKey string
	for key := range input {
		firstKey = key
		break
	}
	return firstKey
}

func GetLastStringKeyFromMap[V any](input map[string]V) string {
	var lastKey string
	inputSize := len(input)
	for key := range input {
		if inputSize == 1 {
			lastKey = key
		}
		inputSize--
	}
	return lastKey
}
