package solidity

import (
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type BoolSerializerTestSuite struct {
	suite.Suite
}

func TestBoolSerializerTestSuite(t *testing.T) {
	suite.Run(t, new(BoolSerializerTestSuite))
}

func (s *BoolSerializerTestSuite) TestDeserializeShouldReturnErrorWhenReceivedInvalidBool() {
	invalidBool := gofakeit.Word()
	boolHandler := NewBoolHandler()
	err := boolHandler.Deserialize(invalidBool)
	assert.ErrorIs(s.T(), err, ErrInvalidBool)
}

func (s *BoolSerializerTestSuite) TestMapStringToBoolShouldReturnErrorWhenReceivedInvalidBool() {
	invalidBool := gofakeit.Word()
	boolHandler := NewBoolHandler()
	err := boolHandler.Deserialize(invalidBool)
	assert.ErrorIs(s.T(), err, ErrInvalidBool)
}

func (s *BoolSerializerTestSuite) TestMapStringToBoolShouldReturnErrorWhenReceivedValidBool() {
	validBoolValues := map[string]bool{
		"true":  true,
		"TRUE":  true,
		"t":     true,
		"T":     true,
		"false": false,
		"FALSE": false,
		"f":     false,
		"F":     false,
	}
	boolHandler := NewBoolHandler()

	for input, expectedReturn := range validBoolValues {
		err := boolHandler.Deserialize(input)
		assert.Nil(s.T(), err)
		assert.Equal(s.T(), expectedReturn, boolHandler.GetValue())
	}
}
