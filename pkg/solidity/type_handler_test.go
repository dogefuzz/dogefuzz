package solidity

import (
	"reflect"
	"testing"

	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/dogefuzz/dogefuzz/test/generators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TypeHandlerTestSuite struct {
	suite.Suite
}

func TestTypeHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(TypeHandlerTestSuite))
}

func (s *TypeHandlerTestSuite) TestGetTypeHandler_ShouldReturnValidBoolHandler_WhenReceiveValidBoolAbiType() {
	boolTyp := generators.BoolTypeGen()
	expectedHandlerType := reflect.TypeOf(NewBoolHandler())

	value, err := GetTypeHandler(boolTyp)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), expectedHandlerType, reflect.TypeOf(value))
}

func (s *TypeHandlerTestSuite) TestGetTypeHandler_ShouldReturnValidUint8Handler_WhenReceiveValidUint8AbiType() {
	uintTyp := generators.UintTypeGen(8)
	expectedHandlerType := reflect.TypeOf(NewUint8Handler())

	value, err := GetTypeHandler(uintTyp)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), expectedHandlerType, reflect.TypeOf(value))
}

func (s *TypeHandlerTestSuite) TestGetTypeHandler_ShouldReturnValidUint16Handler_WhenReceiveValidUint16AbiType() {
	uintTyp := generators.UintTypeGen(16)
	expectedHandlerType := reflect.TypeOf(NewUint16Handler())

	value, err := GetTypeHandler(uintTyp)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), expectedHandlerType, reflect.TypeOf(value))
}

func (s *TypeHandlerTestSuite) TestGetTypeHandler_ShouldReturnValidUint32Handler_WhenReceiveValidUint32AbiType() {
	uintTyp := generators.UintTypeGen(32)
	expectedHandlerType := reflect.TypeOf(NewUint32Handler())

	value, err := GetTypeHandler(uintTyp)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), expectedHandlerType, reflect.TypeOf(value))
}

func (s *TypeHandlerTestSuite) TestGetTypeHandler_ShouldReturnValidUint64Handler_WhenReceiveValidUint64AbiType() {
	uintTyp := generators.UintTypeGen(64)
	expectedHandlerType := reflect.TypeOf(NewUint64Handler())

	value, err := GetTypeHandler(uintTyp)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), expectedHandlerType, reflect.TypeOf(value))
}

func (s *TypeHandlerTestSuite) TestGetTypeHandler_ShouldReturnValidUnsignedBigIntHandler_WhenReceiveValidUintXAbiType() {
	exceptionSet := common.NewSet(8, 16, 32, 64)

	for bitSize := 8; bitSize <= 256; bitSize += 8 {
		if exceptionSet.Has(bitSize) {
			continue
		}

		uintTyp := generators.UintTypeGen(bitSize)
		expectedHandlerType := reflect.TypeOf(NewUnsignedBigIntHandler(bitSize))

		value, err := GetTypeHandler(uintTyp)

		assert.Nil(s.T(), err)
		assert.Equal(s.T(), expectedHandlerType, reflect.TypeOf(value))
	}
}

func (s *TypeHandlerTestSuite) TestGetTypeHandler_ShouldReturnValidInt8Handler_WhenReceiveValidInt8AbiType() {
	intTyp := generators.IntTypeGen(8)
	expectedHandlerType := reflect.TypeOf(NewInt8Handler())

	value, err := GetTypeHandler(intTyp)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), expectedHandlerType, reflect.TypeOf(value))
}

func (s *TypeHandlerTestSuite) TestGetTypeHandler_ShouldReturnValidInt16Handler_WhenReceiveValidInt16AbiType() {
	intTyp := generators.IntTypeGen(16)
	expectedHandlerType := reflect.TypeOf(NewInt16Handler())

	value, err := GetTypeHandler(intTyp)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), expectedHandlerType, reflect.TypeOf(value))
}

func (s *TypeHandlerTestSuite) TestGetTypeHandler_ShouldReturnValidInt32Handler_WhenReceiveValidInt32AbiType() {
	intTyp := generators.IntTypeGen(32)
	expectedHandlerType := reflect.TypeOf(NewInt32Handler())

	value, err := GetTypeHandler(intTyp)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), expectedHandlerType, reflect.TypeOf(value))
}

func (s *TypeHandlerTestSuite) TestGetTypeHandler_ShouldReturnValidInt64Handler_WhenReceiveValidInt64AbiType() {
	intTyp := generators.IntTypeGen(64)
	expectedHandlerType := reflect.TypeOf(NewInt64Handler())

	value, err := GetTypeHandler(intTyp)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), expectedHandlerType, reflect.TypeOf(value))
}

func (s *TypeHandlerTestSuite) TestGetTypeHandler_ShouldReturnValidSignedBigIntHandler_WhenReceiveValidIntXAbiType() {
	exceptionSet := common.NewSet(8, 16, 32, 64)

	for bitSize := 8; bitSize <= 256; bitSize += 8 {
		if exceptionSet.Has(bitSize) {
			continue
		}

		intTyp := generators.IntTypeGen(bitSize)
		expectedHandlerType := reflect.TypeOf(NewSignedBigIntHandler(bitSize))

		value, err := GetTypeHandler(intTyp)

		assert.Nil(s.T(), err)
		assert.Equal(s.T(), expectedHandlerType, reflect.TypeOf(value))
	}
}

func (s *TypeHandlerTestSuite) TestGetTypeHandler_ShouldReturnValidStringHandler_WhenReceiveValidStringAbiType() {
	stringTyp := generators.StringTypeGen()
	expectedHandlerType := reflect.TypeOf(NewStringHandler())

	value, err := GetTypeHandler(stringTyp)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), expectedHandlerType, reflect.TypeOf(value))
}
