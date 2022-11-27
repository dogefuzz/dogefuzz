package vandal

import (
	"regexp"
	"strconv"
	"strings"
)

func readSlicePropertyLine(pattern string, line string) []string {
	r := regexp.MustCompile(pattern + `: \[?([^\[\]]*)\]?`)
	match := r.FindStringSubmatch(line)
	elements := strings.Split(match[1], ",")
	result := make([]string, 0)
	for _, element := range elements {
		value := strings.Trim(element, " ")
		if value != "" {
			result = append(result, value)
		}
	}
	return result
}

func readIntPropertyLine(pattern string, line string) uint64 {
	r := regexp.MustCompile(pattern + `: (.*)`)
	match := r.FindStringSubmatch(line)
	result, _ := strconv.ParseUint(match[1], 10, 64)
	return result
}

func readStringPropertyLine(pattern string, line string) string {
	r := regexp.MustCompile(pattern + `: (.*)`)
	match := r.FindStringSubmatch(line)
	return match[1]
}
