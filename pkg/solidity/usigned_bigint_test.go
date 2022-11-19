package solidity

import (
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/dogefuzz/dogefuzz/test/generators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UnsignedBigIntHandlerTestSuite struct {
	suite.Suite
}

func TestUnsignedBigIntHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(UnsignedBigIntHandlerTestSuite))
}

func (s *UnsignedBigIntHandlerTestSuite) TestDeserialize_ShouldReturnError_WhenReceivedInvalidNumber() {
	invalidNumber := gofakeit.Word()
	exceptionSet := common.NewSet(8, 16, 32, 64)

	for bitSize := 8; bitSize <= 256; bitSize += 8 {
		if exceptionSet.Has(bitSize) {
			continue
		}
		handler := NewUnsignedBigIntHandler(bitSize)
		err := handler.Deserialize(invalidNumber)
		assert.ErrorIs(s.T(), err, ErrInvalidUnsignedBigInt)
	}
}

func (s *UnsignedBigIntHandlerTestSuite) TestDeserialize_ShouldReturnError_WhenReceivedUnderflowedNumber() {
	underflowedNumberAsString := "-1"
	exceptionSet := common.NewSet(8, 16, 32, 64)

	for bitSize := 8; bitSize <= 256; bitSize += 8 {
		if exceptionSet.Has(bitSize) {
			continue
		}
		handler := NewUnsignedBigIntHandler(bitSize)
		err := handler.Deserialize(underflowedNumberAsString)
		assert.ErrorIs(s.T(), err, ErrInvalidUnsignedBigInt)
	}
}

func (s *UnsignedBigIntHandlerTestSuite) TestDeserialize_ShouldReturnError_WhenReceivedOverflowedNumber() {
	exceptionSet := common.NewSet(8, 16, 32, 64)

	for bitSize := 8; bitSize <= 256; bitSize += 8 {
		if exceptionSet.Has(bitSize) {
			continue
		}
		overflowedNumberAsString := generators.OverflowedNumberAsStringGen(bitSize)
		handler := NewUnsignedBigIntHandler(bitSize)

		err := handler.Deserialize(overflowedNumberAsString)

		assert.ErrorIs(s.T(), err, ErrInvalidUnsignedBigInt)
	}
}

func (s *SignedBigIntHandlerTestSuite) TestSequelize_ShouldReturnBigInt_WhenReceivedAValidUnsignedBigInt() {
	rand.Seed(time.Now().Unix())
	exceptionSet := common.NewSet(8, 16, 32, 64)

	for bitSize := 8; bitSize <= 256; bitSize += 8 {
		if exceptionSet.Has(bitSize) {
			continue
		}

		validUnsignedBigInt := generators.UnsignedBigIntGen(bitSize)
		validUnsignedBigIntAsString := validUnsignedBigInt.String()
		handler := NewUnsignedBigIntHandler(bitSize)

		err := handler.Deserialize(validUnsignedBigIntAsString)

		assert.Nil(s.T(), err)
		assert.True(s.T(), validUnsignedBigInt.Cmp(handler.GetValue().(*big.Int)) == 0)
	}
}
