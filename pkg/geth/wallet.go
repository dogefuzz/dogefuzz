package geth

import (
	"crypto/ecdsa"
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

var ErrInvalidPublicKey = errors.New("could not parse the private key from the provided hex string")
var ErrCouldNotDerivePublicKey = errors.New("could not derive public key from private key")

type wallet struct {
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
}

func NewWallet() (*wallet, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}

	publicKey, ok := privateKey.Public().(*ecdsa.PublicKey)
	if !ok {
		return nil, ErrCouldNotDerivePublicKey
	}

	return &wallet{privateKey, publicKey}, nil
}

func NewWalletFromPrivateKeyHex(hex string) (*wallet, error) {
	privateKey, err := crypto.HexToECDSA(hex)
	if err != nil {
		return nil, ErrInvalidPublicKey
	}

	publicKey, ok := privateKey.Public().(*ecdsa.PublicKey)
	if !ok {
		return nil, ErrCouldNotDerivePublicKey
	}

	return &wallet{privateKey, publicKey}, nil
}

func (w *wallet) GetPrivateKey() *ecdsa.PrivateKey {
	return w.privateKey
}

func (w *wallet) GetPrivateKeyHex() string {
	privateKeyAsBytes := crypto.FromECDSA(w.privateKey)
	return hexutil.Encode(privateKeyAsBytes)[2:] // Removing 0x prefix
}

func (w *wallet) GetPublicKey() *ecdsa.PublicKey {
	return w.publicKey
}

func (w *wallet) GetPublicKeyHex() string {
	publicKeyAsBytes := crypto.FromECDSAPub(w.publicKey)
	return hexutil.Encode(publicKeyAsBytes)[4:] // Removing 0x04 prefix
}

func (w *wallet) GetAddress() common.Address {
	return crypto.PubkeyToAddress(*w.publicKey)
}
