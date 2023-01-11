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

type Uint64HandlerTestSuite struct {
	suite.Suite
}

func TestUint64HandlerTestSuite(t *testing.T) {
	suite.Run(t, new(Uint64HandlerTestSuite))
}

func (s *Uint64HandlerTestSuite) TestDeserialize_ShouldReturnError_WhenReceivedInvalidNumber() {
	invalidNumber := gofakeit.Word()
	handler := NewUint64Handler()

	err := handler.Deserialize(invalidNumber)

	assert.ErrorIs(s.T(), err, ErrInvalidUint64)
}

func (s *Uint64HandlerTestSuite) TestDeserialize_ShouldReturnError_WhenReceivedUnderflowedNumber() {
	underflowedNumberAsString := "-1"
	handler := NewUint64Handler()

	err := handler.Deserialize(underflowedNumberAsString)

	assert.ErrorIs(s.T(), err, ErrInvalidUint64)
}

func (s *Uint64HandlerTestSuite) TestDeserialize_ShouldReturnError_WhenReceivedOverflowedNumber() {
	overflowedNumberAsString := generators.OverflowedNumberAsStringGen(64)
	handler := NewUint64Handler()

	err := handler.Deserialize(overflowedNumberAsString)

	assert.ErrorIs(s.T(), err, ErrInvalidUint64)
}

func (s *Uint64HandlerTestSuite) TestDeserialize_ShouldFillValueWithValidUint64_WhenReceivedAValidUint64() {
	rand.Seed(common.Now().Unix())
	validUint64 := rand.Uint64()
	validUint64AsString := strconv.FormatUint(validUint64, 10)
	handler := NewUint64Handler()

	err := handler.Deserialize(validUint64AsString)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), uint64(validUint64), handler.GetValue())
}
