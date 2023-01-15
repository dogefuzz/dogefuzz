package common

import (
	"math/rand"
	"sort"
	"strings"
	"time"
)

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
	set := make(map[string]bool)
	for _, elem := range a {
		set[elem] = true
	}

	for _, elem := range b {
		if _, ok := set[elem]; !ok {
			set[elem] = true
		}
	}

	result := make([]string, 0, len(set))
	for elem := range set {
		result = append(result, elem)
	}

	sort.Strings(result)
	return result
}

func RandomSubSlice[T any](s []T) []T {
	if len(s) == 0 {
		return make([]T, 0)
	}
	rand.Seed(time.Now().UnixNano())
	subSliceLength := rand.Intn(len(s))
	subSlice := make([]T, subSliceLength)

	leftIdx := 0
	currentLength := subSliceLength
	for i := 0; i < subSliceLength; i++ {
		idx := rand.Intn(currentLength) + leftIdx
		subSlice[i] = s[idx]
		Swap(s, idx, leftIdx)
		currentLength--
		leftIdx++
	}

	return subSlice
}

func Swap[T any](s []T, i, j int) {
	temp := s[i]
	s[i] = s[j]
	s[j] = temp
}

func Contains[T comparable](s []T, target T) bool {
	for _, element := range s {
		if element == target {
			return true
		}
	}
	return false
}
