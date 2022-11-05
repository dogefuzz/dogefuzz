package solc

import (
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CommonIntegrationTestSuite struct {
	suite.Suite
}

func TestCommonIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(CommonIntegrationTestSuite))
}

func (s *CommonIntegrationTestSuite) TestGetSolidityVersions() {
	versions, err := getDescendingOrderedVersionsFromSolidyBinariesEndpoint()
	assert.Nil(s.T(), err)
	assert.NotEmpty(s.T(), versions)
	for _, version := range versions {
		_, err := semver.NewVersion(version)
		assert.Nil(s.T(), err)
	}
}
