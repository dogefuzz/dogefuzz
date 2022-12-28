package common

import "strings"

func GetUniqueSlice[T comparable](slice []T) []T {
	keys := make(map[T]bool)
	list := []T{}
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func JoinOracleTypeList(types []OracleType) string {
	words := make([]string, 0)
	for _, t := range types {
		words = append(words, string(t))
	}
	return strings.Join(words, ";")
}

func SplitOracleTypeString(input string) []OracleType {
	words := strings.Split(input, ";")
	result := make([]OracleType, 0)
	for _, w := range words {
		result = append(result, OracleType(w))
	}
	return result
}

func MergeSortedSlices(a []string, b []string) []string {
	// TODO: merge two sorted array and remove duplicated elements
	return nil
}
