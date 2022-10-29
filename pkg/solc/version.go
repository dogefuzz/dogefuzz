package solc

import (
	"regexp"
	"strconv"
)

type Version struct {
	Patch int
	Minor int
	Major int
}

func FromString(version string) *Version {
	r, _ := regexp.Compile(`([0-9]+)\.([0-9]+)\.([0-9]+)`)
	matches := r.FindStringSubmatch(version)
	patch, _ := strconv.Atoi(matches[3])
	minor, _ := strconv.Atoi(matches[2])
	major, _ := strconv.Atoi(matches[1])

	return &Version{
		Major: major,
		Minor: minor,
		Patch: patch,
	}
}
