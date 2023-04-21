package common

func SumOccurrencesOfStringList(targetGroup []string, source []string) uint64 {
	var count uint64 = 0
	for _, element := range source {
		for _, target := range targetGroup {
			if element == target {
				count++
			}
		}
	}
	return count
}
