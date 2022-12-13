package geth

import (
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type WalletTestSuite struct {
	suite.Suite
}

func TestWalletTestSuite(t *testing.T) {
	suite.Run(t, new(WalletTestSuite))
}

func (s *WalletTestSuite) TestNewWallet_ShouldGenerateUniqueWallet() {
	privateKeysSet := common.NewSet[string]()
	errorsList := make([]error, 0)

	for i := 0; i < 100; i++ {
		wallet, err := NewWallet()
		privateKeysSet.Add(wallet.GetPrivateKeyHex())

		if err != nil {
			errorsList = append(errorsList, err)
		}
	}

	assert.Equal(s.T(), privateKeysSet.Len(), 100)
	assert.Empty(s.T(), errorsList)
}

func (s *WalletTestSuite) TestNewWalletFromPrivateKeyHex_ShouldGenerateEquivalentWallet_WhenReceivedAWalletPrivateKeyHex() {
	wallet, _ := NewWallet()
	walletByPrivateKeyHex, err := NewWalletFromPrivateKeyHex(wallet.GetPrivateKeyHex())

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), wallet.GetPublicKeyHex(), walletByPrivateKeyHex.GetPublicKeyHex())
}

func (s *WalletTestSuite) TestNewWalletFromPrivateKeyHex_ShouldReturnErrCouldNotDerivePublicKey_WhenReceivedAnInvalidWalletPrivateKeyHex() {
	invalidPrivateKeyHex := gofakeit.HexUint256()
	_, err := NewWalletFromPrivateKeyHex(invalidPrivateKeyHex)

	assert.ErrorIs(s.T(), err, ErrInvalidPublicKey)
}
