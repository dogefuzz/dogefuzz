package geth

import (
	"testing"

	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type WalletIntegrationTestSuite struct {
	suite.Suite
}

func TestWalletIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(WalletIntegrationTestSuite))
}

func (s *WalletIntegrationTestSuite) TestNewWallet_ShouldGenerateUniqueWallet() {
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

func (s *WalletIntegrationTestSuite) TestNewWalletFromPrivateKeyHex_ShouldGenerateEquivalentWallet_WhenReceivedAWalletPrivateKeyHex() {
	wallet, _ := NewWallet()
	walletByPrivateKeyHex, err := NewWalletFromPrivateKeyHex(wallet.GetPrivateKeyHex())

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), wallet.GetPublicKeyHex(), walletByPrivateKeyHex.GetPublicKeyHex())
}
