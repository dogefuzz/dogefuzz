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

type SignedBigIntHandlerTestSuite struct {
	suite.Suite
}

func TestSignedBigIntHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(SignedBigIntHandlerTestSuite))
}

func (s *SignedBigIntHandlerTestSuite) TestDeserialize_ShouldReturnError_WhenReceivedInvalidNumber() {
	invalidNumber := gofakeit.Word()
	exceptionSet := common.NewSet(8, 16, 32, 64)

	for bitSize := 8; bitSize <= 256; bitSize += 8 {
		if exceptionSet.Has(bitSize) {
			continue
		}
		handler := NewSignedBigIntHandler(bitSize)
		err := handler.Deserialize(invalidNumber)
		assert.ErrorIs(s.T(), err, ErrInvalidSignedBigInt)
	}
}

func (s *SignedBigIntHandlerTestSuite) TestDeserialize_ShouldReturnError_WhenReceivedUnderflowedNumber() {
	exceptionSet := common.NewSet(8, 16, 32, 64)

	for bitSize := 8; bitSize <= 256; bitSize += 8 {
		if exceptionSet.Has(bitSize) {
			continue
		}
		underflowedNumberAsString := generators.UnderflowedNumberAsStringGen(bitSize)
		handler := NewSignedBigIntHandler(bitSize)
		err := handler.Deserialize(underflowedNumberAsString)
		assert.ErrorIs(s.T(), err, ErrInvalidSignedBigInt)
	}
}

func (s *SignedBigIntHandlerTestSuite) TestDeserialize_ShouldReturnError_WhenReceivedOverflowedNumber() {
	exceptionSet := common.NewSet(8, 16, 32, 64)

	for bitSize := 8; bitSize <= 256; bitSize += 8 {
		if exceptionSet.Has(bitSize) {
			continue
		}
		overflowedNumberAsString := generators.OverflowedNumberAsStringGen(bitSize)
		handler := NewSignedBigIntHandler(bitSize)

		err := handler.Deserialize(overflowedNumberAsString)

		assert.ErrorIs(s.T(), err, ErrInvalidSignedBigInt)
	}
}

func (s *SignedBigIntHandlerTestSuite) TestSequelize_ShouldReturnBigInt_WhenReceivedAValidSignedBigInt() {
	rand.Seed(time.Now().Unix())
	exceptionSet := common.NewSet(8, 16, 32, 64)

	for bitSize := 8; bitSize <= 256; bitSize += 8 {
		if exceptionSet.Has(bitSize) {
			continue
		}

		validSignedBigInt := generators.SignedBigIntGen(bitSize)
		validSignedBigIntAsString := validSignedBigInt.String()
		handler := NewSignedBigIntHandler(bitSize)

		err := handler.Deserialize(validSignedBigIntAsString)

		assert.Nil(s.T(), err)
		assert.True(s.T(), validSignedBigInt.Cmp(handler.GetValue().(*big.Int)) == 0)
	}
}
