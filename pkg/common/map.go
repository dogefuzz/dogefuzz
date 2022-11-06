package common

func GetFirstStringKeyFromMap[V any](input map[string]V) string {
	var firstKey string
	for key := range input {
		firstKey = key
		break
	}
	return firstKey
}
