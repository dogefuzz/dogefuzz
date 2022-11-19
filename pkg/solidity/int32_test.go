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

type Int32HandlerTestSuite struct {
	suite.Suite
}

func TestInt32HandlerTestSuite(t *testing.T) {
	suite.Run(t, new(Int32HandlerTestSuite))
}

func (s *Int32HandlerTestSuite) TestDeserialize_ShouldReturnError_WhenReceivedInvalidNumber() {
	invalidNumber := gofakeit.Word()
	handler := NewInt32Handler()

	err := handler.Deserialize(invalidNumber)

	assert.ErrorIs(s.T(), err, ErrInvalidInt32)
}

func (s *Int32HandlerTestSuite) TestDeserialize_ShouldReturnError_WhenReceivedUnderflowedNumber() {
	underflowedNumberAsString := generators.UnderflowedNumberAsStringGen(32)
	handler := NewInt32Handler()

	err := handler.Deserialize(underflowedNumberAsString)

	assert.ErrorIs(s.T(), err, ErrInvalidInt32)
}

func (s *Int32HandlerTestSuite) TestDeserialize_ShouldReturnError_WhenReceivedOverflowedNumber() {
	overflowedNumberAsString := generators.OverflowedNumberAsStringGen(32)
	handler := NewInt32Handler()

	err := handler.Deserialize(overflowedNumberAsString)

	assert.ErrorIs(s.T(), err, ErrInvalidInt32)
}

func (s *Int32HandlerTestSuite) TestDeserialize_ShouldReturnInt32_WhenReceivedAValidInt32() {
	rand.Seed(time.Now().Unix())
	validInt32 := int32(generators.SignedBigIntGen(32).Int64())
	validInt32AsString := strconv.Itoa(int(validInt32))
	handler := NewInt32Handler()

	err := handler.Deserialize(validInt32AsString)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), validInt32, handler.GetValue())
}
