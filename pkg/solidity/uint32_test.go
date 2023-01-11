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

type Uint32HandlerTestSuite struct {
	suite.Suite
}

func TestUint32HandlerTestSuite(t *testing.T) {
	suite.Run(t, new(Uint32HandlerTestSuite))
}

func (s *Uint32HandlerTestSuite) TestDeserialize_ShouldReturnError_WhenReceivedInvalidNumber() {
	invalidNumber := gofakeit.Word()
	handler := NewUint32Handler()

	err := handler.Deserialize(invalidNumber)

	assert.ErrorIs(s.T(), err, ErrInvalidUint32)
}

func (s *Uint32HandlerTestSuite) TestDeserialize_ShouldReturnError_WhenReceivedUnderflowedNumber() {
	underflowedNumberAsString := "-1"
	handler := NewUint32Handler()

	err := handler.Deserialize(underflowedNumberAsString)

	assert.ErrorIs(s.T(), err, ErrInvalidUint32)
}

func (s *Uint32HandlerTestSuite) TestDeserialize_ShouldReturnError_WhenReceivedOverflowedNumber() {
	overflowedNumberAsString := generators.OverflowedNumberAsStringGen(32)
	handler := NewUint32Handler()

	err := handler.Deserialize(overflowedNumberAsString)

	assert.ErrorIs(s.T(), err, ErrInvalidUint32)
}

func (s *Uint32HandlerTestSuite) TestDeserialize_ShouldFillValueWithValidUint32_WhenReceivedAValidUint32() {
	rand.Seed(common.Now().Unix())
	validUint32 := rand.Uint32()
	validUint32AsString := strconv.FormatUint(uint64(validUint32), 10)
	handler := NewUint32Handler()

	err := handler.Deserialize(validUint32AsString)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), uint32(validUint32), handler.GetValue())
}
