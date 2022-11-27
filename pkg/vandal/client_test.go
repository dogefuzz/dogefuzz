package vandal

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type VandalIntegrationTestSuite struct {
	suite.Suite
}

func TestVandalIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(VandalIntegrationTestSuite))
}

func (s *VandalIntegrationTestSuite) TestDecompile_ReturnNil() {
	// c := NewVandalClient()
	// blocks, functions, err := c.Decompile()
	// assert.Equal(s.T(), len(blocks), 0)
	// assert.Equal(s.T(), len(functions), 0)
	// assert.NotNil(s.T(), err)
}
