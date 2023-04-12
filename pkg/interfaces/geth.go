package interfaces

import (
	"context"
	"crypto/ecdsa"
	"math/big"

	"github.com/dogefuzz/dogefuzz/pkg/common"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

type Agent interface {
	Send(ctx context.Context, wallet Wallet, contract *common.Contract, functionName string, value *big.Int, args ...interface{}) (string, error)
	Transfer(ctx context.Context, wallet Wallet, contract *common.Contract, value *big.Int) (string, error)
}

type Deployer interface {
	Deploy(ctx context.Context, contract *common.Contract, args ...interface{}) (string, string, error)
}

type Wallet interface {
	GetPrivateKey() *ecdsa.PrivateKey
	GetPrivateKeyHex() string
	GetPublicKey() *ecdsa.PublicKey
	GetPublicKeyHex() string
	GetAddress() gethcommon.Address
}
