package solidity

import (
	"math"
	"math/rand"
	"strconv"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/dogefuzz/dogefuzz/test/generators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type Uint8HandlerTestSuite struct {
	suite.Suite
}

func TestUint8HandlerTestSuite(t *testing.T) {
	suite.Run(t, new(Uint8HandlerTestSuite))
}

func (s *Uint8HandlerTestSuite) TestDeserialize_ShouldReturnError_WhenReceivedInvalidNumber() {
	invalidNumber := gofakeit.Word()
	handler := NewUint8Handler()

	err := handler.Deserialize(invalidNumber)

	assert.ErrorIs(s.T(), err, ErrInvalidUint8)
}

func (s *Uint8HandlerTestSuite) TestDeserialize_ShouldReturnError_WhenReceivedUnderflowedNumber() {
	underflowedNumberAsString := "-1"
	handler := NewUint8Handler()

	err := handler.Deserialize(underflowedNumberAsString)

	assert.ErrorIs(s.T(), err, ErrInvalidUint8)
}

func (s *Uint8HandlerTestSuite) TestDeserialize_ShouldReturnError_WhenReceivedOverflowedNumber() {
	overflowedNumberAsString := generators.OverflowedNumberAsStringGen(8)
	handler := NewUint8Handler()

	err := handler.Deserialize(overflowedNumberAsString)

	assert.ErrorIs(s.T(), err, ErrInvalidUint8)
}

func (s *Uint8HandlerTestSuite) TestDeserialize_ShouldFillValueWithValidUint8_WhenReceivedAValidUint8() {
	rand.Seed(common.Now().Unix())
	validUint8 := rand.Intn(math.MaxUint8)
	validUint8AsString := strconv.Itoa(validUint8)
	handler := NewUint8Handler()

	err := handler.Deserialize(validUint8AsString)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), uint8(validUint8), handler.GetValue())
}
