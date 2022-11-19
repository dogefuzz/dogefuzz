package solidity

import (
	"math"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/dogefuzz/dogefuzz/test/generators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type Uint16HandlerTestSuite struct {
	suite.Suite
}

func TestUint16HandlerTestSuite(t *testing.T) {
	suite.Run(t, new(Uint16HandlerTestSuite))
}

func (s *Uint16HandlerTestSuite) TestDeserialize_ShouldReturnError_WhenReceivedInvalidNumber() {
	invalidNumber := gofakeit.Word()
	handler := NewUint16Handler()

	err := handler.Deserialize(invalidNumber)

	assert.ErrorIs(s.T(), err, ErrInvalidUint16)
}

func (s *Uint16HandlerTestSuite) TestDeserialize_ShouldReturnError_WhenReceivedUnderflowedNumber() {
	underflowedNumberAsString := generators.UnderflowedNumberAsStringGen(16)
	handler := NewUint16Handler()

	err := handler.Deserialize(underflowedNumberAsString)

	assert.ErrorIs(s.T(), err, ErrInvalidUint16)
}

func (s *Uint16HandlerTestSuite) TestDeserialize_ShouldReturnError_WhenReceivedOverflowedNumber() {
	overflowedNumberAsString := generators.OverflowedNumberAsStringGen(16)
	handler := NewUint16Handler()

	err := handler.Deserialize(overflowedNumberAsString)

	assert.ErrorIs(s.T(), err, ErrInvalidUint16)
}

func (s *Uint16HandlerTestSuite) TestDeserialize_ShouldReturnUint16_WhenReceivedAValidUint16() {
	rand.Seed(time.Now().Unix())
	validUint16 := rand.Intn(math.MaxUint16)
	validUint16AsString := strconv.Itoa(validUint16)
	handler := NewUint16Handler()

	err := handler.Deserialize(validUint16AsString)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), uint16(validUint16), handler.GetValue())
}
