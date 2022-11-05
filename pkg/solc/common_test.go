package solc

import (
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CommonTestSuite struct {
	suite.Suite
}

func TestCommonTestSuite(t *testing.T) {
	suite.Run(t, new(CommonTestSuite))
}

func (s *CommonTestSuite) TestgetMaxVersionBasedOnContraintWhenAnInvalidConstraintIsPassed() {
	invalidConstraint := "~1.2.0"
	versions := []string{"1.3.0", "1.2.9", "1.2.0"}

	version, err := getMaxVersionBasedOnContraint(versions, invalidConstraint)
	assert.Nil(s.T(), err)
	expectedVersion, _ := semver.NewVersion("1.2.9")
	assert.True(s.T(), version.Equal(expectedVersion))
}

func (s *CommonTestSuite) TestgetMaxVersionBasedOnContraintWhenReceivedATildeRangeComparison() {
	invalidConstraint := "~1.2.0"
	versions := []string{"1.3.0", "1.2.9", "1.2.0"}

	version, err := getMaxVersionBasedOnContraint(versions, invalidConstraint)
	assert.Nil(s.T(), err)
	expectedVersion, _ := semver.NewVersion("1.2.9")
	assert.True(s.T(), version.Equal(expectedVersion))
}

func (s *CommonTestSuite) TestgetMaxVersionBasedOnContraintWhenReceivedACaretRangeComparison() {
	invalidConstraint := "^1.2.0"
	versions := []string{"2.0.0", "1.9.0", "1.2.0"}

	version, err := getMaxVersionBasedOnContraint(versions, invalidConstraint)
	assert.Nil(s.T(), err)
	expectedVersion, _ := semver.NewVersion("1.9.0")
	assert.True(s.T(), version.Equal(expectedVersion))
}

func (s *CommonTestSuite) TestgetMaxVersionBasedOnContraintWhenReceivedARangeComparison() {
	invalidConstraint := ">=1.2.0 <= 1.4.0"
	versions := []string{"1.4.1", "1.4.0", "1.3.0", "1.2.0"}

	version, err := getMaxVersionBasedOnContraint(versions, invalidConstraint)
	assert.Nil(s.T(), err)
	expectedVersion, _ := semver.NewVersion("1.4.0")
	assert.True(s.T(), version.Equal(expectedVersion))
}

func (s *CommonTestSuite) TestgetMaxVersionBasedOnContraintWhenInvalidInputIsProvided() {
	invalidConstraint := "invalid-constraint"
	versions := []string{"1.2.0"}

	version, err := getMaxVersionBasedOnContraint(versions, invalidConstraint)
	assert.Nil(s.T(), version)
	assert.EqualError(s.T(), err, ErrInvalidConstraint.Error())
}

func (s *CommonTestSuite) TestextractVersionConstraintFromSourceWhenAWellDefinedFileIsProvided() {
	const WELL_DEFINED_CONTRACT = `
	pragma solidity ^0.8.13;

	contract HelloWorld {
		string public greet = "Hello World!";
	}`

	constraint, _ := extractVersionConstraintFromSource(WELL_DEFINED_CONTRACT)

	assert.Equal(s.T(), "^0.8.13", constraint)
}

func (s *CommonTestSuite) TestExtractVersionWhenTheSolidityFileStartsWithAComment() {
	const CONTRACT_WITH_CONTENT = `
	// SPDX-License-Identifier: MIT
	pragma solidity ^0.8.13;

	contract HelloWorld {
		string public greet = "Hello World!";
	}
	`

	constraint, _ := extractVersionConstraintFromSource(CONTRACT_WITH_CONTENT)

	assert.Equal(s.T(), "^0.8.13", constraint)
}

func (s *CommonTestSuite) TestExtractVersionWhenTheSolidityFileContainsTheVersionMoreThanOnce() {
	const CONTRACT_WITH_MORE_THAN_ONE_VERSION_DEFINITION = `
	// SPDX-License-Identifier: MIT
	pragma solidity ^0.8.13;
	//pragma solidity ^0.10.13;
	//pragma solidity ^0.9.13;
	// pragma solidity ^0.11.13;
	// pragma solidity ^0.12.13;

	contract HelloWorld {
		string public greet = "Hello World!";
	}
	`

	constraint, _ := extractVersionConstraintFromSource(CONTRACT_WITH_MORE_THAN_ONE_VERSION_DEFINITION)

	assert.Equal(s.T(), "^0.8.13", constraint)
}

func (s *CommonTestSuite) TestExtractVersionWhenThe() {
	const CONTRACT_WITH_MORE_THAN_ONE_VERSION_DEFINITION = `
	// SPDX-License-Identifier: MIT
	//pragma solidity ^0.10.13;
	//pragma solidity ^0.9.13;
	// pragma solidity ^0.11.13;
	// pragma solidity ^0.12.13;
	pragma solidity ^0.8.13;

	contract HelloWorld {
		string public greet = "Hello World!";
	}
	`

	constraint, _ := extractVersionConstraintFromSource(CONTRACT_WITH_MORE_THAN_ONE_VERSION_DEFINITION)

	assert.Equal(s.T(), "^0.8.13", constraint)
}

func (s *CommonTestSuite) TestextractVersionConstraintFromSource() {
	const CONTRACT_WITH_CONTENT = `
	// SPDX-License-Identifier: MIT
	pragma solidity ^0.8.13;

	contract HelloWorld {
		string public greet = "Hello World!";
	}
	`

	constraint, _ := extractVersionConstraintFromSource(CONTRACT_WITH_CONTENT)

	assert.Equal(s.T(), "^0.8.13", constraint)
}

func (s *CommonTestSuite) TestextractVersionConstraintFromSourceWhenThereIsARange() {
	const CONTRACT_WITH_RANGE = `
	// SPDX-License-Identifier: MIT
	pragma solidity >=0.8.13 <= 0.8.14;

	contract HelloWorld {
		string public greet = "Hello World!";
	}
	`

	constraint, _ := extractVersionConstraintFromSource(CONTRACT_WITH_RANGE)

	assert.Equal(s.T(), ">=0.8.13 <= 0.8.14", constraint)
}

func (s *CommonTestSuite) TestextractVersionConstraintFromSourceWhenTheContractIsEmpty() {
	const EMPTY_CONTRACT = ""

	constraint, err := extractVersionConstraintFromSource(EMPTY_CONTRACT)

	assert.Empty(s.T(), constraint)
	assert.EqualError(s.T(), err, ErrVersionNotFound.Error())
}
