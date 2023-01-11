package solidity

import (
	"math/rand"
	"strconv"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/dogefuzz/dogefuzz/test/generators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type Int64HandlerTestSuite struct {
	suite.Suite
}

func TestInt64HandlerTestSuite(t *testing.T) {
	suite.Run(t, new(Int64HandlerTestSuite))
}

func (s *Int64HandlerTestSuite) TestDeserialize_ShouldReturnError_WhenReceivedInvalidNumber() {
	invalidNumber := gofakeit.Word()
	handler := NewInt64Handler()

	err := handler.Deserialize(invalidNumber)

	assert.ErrorIs(s.T(), err, ErrInvalidInt64)
}

func (s *Int64HandlerTestSuite) TestDeserialize_ShouldReturnError_WhenReceivedUnderflowedNumber() {
	underflowedNumberAsString := generators.UnderflowedNumberAsStringGen(64)
	handler := NewInt64Handler()

	err := handler.Deserialize(underflowedNumberAsString)

	assert.ErrorIs(s.T(), err, ErrInvalidInt64)
}

func (s *Int64HandlerTestSuite) TestDeserialize_ShouldReturnError_WhenReceivedOverflowedNumber() {
	overflowedNumberAsString := generators.OverflowedNumberAsStringGen(64)
	handler := NewInt64Handler()

	err := handler.Deserialize(overflowedNumberAsString)

	assert.ErrorIs(s.T(), err, ErrInvalidInt64)
}

func (s *Int64HandlerTestSuite) TestDeserialize_ShouldReturnInt64_WhenReceivedAValidInt64() {
	rand.Seed(common.Now().Unix())
	validInt64 := generators.SignedBigIntGen(64).Int64()
	validInt64AsString := strconv.FormatInt(validInt64, 10)
	handler := NewInt64Handler()

	err := handler.Deserialize(validInt64AsString)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), validInt64, handler.GetValue())
}
