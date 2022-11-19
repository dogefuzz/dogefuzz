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

type Int16HandlerTestSuite struct {
	suite.Suite
}

func TestInt16HandlerTestSuite(t *testing.T) {
	suite.Run(t, new(Int16HandlerTestSuite))
}

func (s *Int16HandlerTestSuite) TestDeserialize_ShouldReturnError_WhenReceivedInvalidNumber() {
	invalidNumber := gofakeit.Word()
	handler := NewInt16Handler()

	err := handler.Deserialize(invalidNumber)

	assert.ErrorIs(s.T(), err, ErrInvalidInt16)
}

func (s *Int16HandlerTestSuite) TestDeserialize_ShouldReturnError_WhenReceivedUnderflowedNumber() {
	underflowedNumberAsString := generators.UnderflowedNumberAsStringGen(16)
	handler := NewInt16Handler()

	err := handler.Deserialize(underflowedNumberAsString)

	assert.ErrorIs(s.T(), err, ErrInvalidInt16)
}

func (s *Int16HandlerTestSuite) TestDeserialize_ShouldReturnError_WhenReceivedOverflowedNumber() {
	overflowedNumberAsString := generators.OverflowedNumberAsStringGen(16)
	handler := NewInt16Handler()

	err := handler.Deserialize(overflowedNumberAsString)

	assert.ErrorIs(s.T(), err, ErrInvalidInt16)
}

func (s *Int16HandlerTestSuite) TestDeserialize_ShouldReturnInt16_WhenReceivedAValidInt16() {
	rand.Seed(time.Now().Unix())
	validInt16 := int16(generators.SignedBigIntGen(16).Int64())
	validInt16AsString := strconv.Itoa(int(validInt16))
	handler := NewInt16Handler()

	err := handler.Deserialize(validInt16AsString)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), validInt16, handler.GetValue())
}
