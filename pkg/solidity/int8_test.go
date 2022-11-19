package solidity

import (
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/dogefuzz/dogefuzz/test/generators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type Int8HandlerTestSuite struct {
	suite.Suite
}

func TestInt8HandlerTestSuite(t *testing.T) {
	suite.Run(t, new(Int8HandlerTestSuite))
}

func (s *Int8HandlerTestSuite) TestDeserialize_ShouldReturnError_WhenReceivedInvalidNumber() {
	invalidNumber := gofakeit.Word()
	handler := NewInt8Handler()

	err := handler.Deserialize(invalidNumber)

	assert.ErrorIs(s.T(), err, ErrInvalidInt8)
}

func (s *Int8HandlerTestSuite) TestDeserialize_ShouldReturnError_WhenReceivedUnderflowedNumber() {
	underflowedNumberAsString := generators.UnderflowedNumberAsStringGen(8)
	handler := NewInt8Handler()

	err := handler.Deserialize(underflowedNumberAsString)

	assert.ErrorIs(s.T(), err, ErrInvalidInt8)
}

func (s *Int8HandlerTestSuite) TestDeserialize_ShouldReturnError_WhenReceivedOverflowedNumber() {
	overflowedNumberAsString := generators.OverflowedNumberAsStringGen(8)
	handler := NewInt8Handler()

	err := handler.Deserialize(overflowedNumberAsString)

	assert.ErrorIs(s.T(), err, ErrInvalidInt8)
}

func (s *Int8HandlerTestSuite) TestDeserialize_ShouldFillValueWithValidInt8_WhenReceivedAValidInt8() {
	rand.Seed(time.Now().Unix())
	validInt8 := int8(generators.SignedBigIntGen(8).Int64())
	validInt8AsString := strconv.Itoa(int(validInt8))
	handler := NewInt8Handler()

	err := handler.Deserialize(validInt8AsString)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), validInt8, handler.GetValue())
}
