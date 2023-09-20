package solc

import (
	"bufio"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/Masterminds/semver/v3"
)

var ErrVersionNotFound = errors.New("the pragma version string was not found")
var ErrNoVersionMatch = errors.New("no solc version matched with the solidity file's constraint")
var ErrInvalidConstraint = errors.New("a invalid constraint string was provided")

func ErrSolidityBinariesListCouldNotBeDownloaded(code string) error {
	return errors.New(fmt.Sprintf("the solidity binaries list could not be downloaded externally (HTTP %s)", code))
}

func getMaxVersionBasedOnContraint(descendingSortedVersions []string, constraintStr string) (*semver.Version, error) {
	constraints, err := semver.NewConstraint(constraintStr)
	if err != nil {
		return nil, ErrInvalidConstraint
	}

	for _, version := range descendingSortedVersions {
		v, err := semver.NewVersion(version)
		if err != nil {
			continue
		}

		if constraints.Check(v) {
			return v, nil
		}
	}

	return nil, ErrNoVersionMatch
}

func extractVersionConstraintFromSource(source string) (string, error) {
	pragmaVersionRegex, _ := regexp.Compile(`^\s*pragma solidity (.*);`)
	lines := strings.Split(source, "\n")
	for _, line := range lines {
		if !pragmaVersionRegex.MatchString(line) {
			continue
		}
		matches := pragmaVersionRegex.FindStringSubmatch(line)
		return matches[1], nil

	}
	return "", ErrVersionNotFound
}

func getDescendingOrderedVersionsFromSolidyBinariesEndpoint() ([]string, error) {
	const SOLIDITY_BINARIES_LIST_ENDPOINT = "https://binaries.soliditylang.org/bin/list.txt"
	resp, err := http.Get(SOLIDITY_BINARIES_LIST_ENDPOINT)
	if err != nil {
		return nil, ErrSolidityBinariesListCouldNotBeDownloaded(err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, ErrSolidityBinariesListCouldNotBeDownloaded(strconv.Itoa(resp.StatusCode))
	}

	var versions []string
	scanner := bufio.NewScanner(resp.Body)
	pattern := regexp.MustCompile("soljson-")
	for scanner.Scan() {
		extratedVersion := pattern.ReplaceAllString(scanner.Text(), "")
		versions = append(versions, extratedVersion)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if sort.StringsAreSorted(versions) {
		sort.Strings(versions)
	}
	return versions, nil
}
