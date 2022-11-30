package common

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type InterfaceTestSuite struct {
	suite.Suite
}

func TestInterfaceTestSuite(t *testing.T) {
	suite.Run(t, new(InterfaceTestSuite))
}

func (s *InterfaceTestSuite) TestConvertStringArrayToInterfaceArray_ShouldReturnAnEmptyArray_WhenProvideAnEmptyArray() {
	args := make([]string, 0)

	result := ConvertStringArrayToInterfaceArray(args)

	assert.Len(s.T(), result, 0)
}

func (s *InterfaceTestSuite) TestConvertStringArrayToInterfaceArray_ShouldReturnANonEmptyArrayOfInterfaces_WhenProvideANonEmptyArrayOfStrings() {
	args := make([]string, 10)

	result := ConvertStringArrayToInterfaceArray(args)

	assert.Len(s.T(), result, len(args))
	for _, resultArg := range result {
		assert.IsType(s.T(), "", resultArg)
	}
}
