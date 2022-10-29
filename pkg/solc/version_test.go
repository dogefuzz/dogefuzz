package solc

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type VersionTestSuite struct {
	suite.Suite
}

func TestVersionTestSuite(t *testing.T) {
	suite.Run(t, new(VersionTestSuite))
}

func (s *VersionTestSuite) TestFromString() {
	version1 := FromString("0.1.0")
	assert.Equal(s.T(), Version{Major: 0, Minor: 1, Patch: 0}, *version1)
	version2 := FromString("1.0.0")
	assert.Equal(s.T(), Version{Major: 1, Minor: 0, Patch: 0}, *version2)
	version3 := FromString("10.0.0")
	assert.Equal(s.T(), Version{Major: 10, Minor: 0, Patch: 0}, *version3)
	version4 := FromString("10.2.0")
	assert.Equal(s.T(), Version{Major: 10, Minor: 2, Patch: 0}, *version4)
	version5 := FromString("0.0.0")
	assert.Equal(s.T(), Version{Major: 0, Minor: 0, Patch: 0}, *version5)
}
