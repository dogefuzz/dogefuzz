package interfaces

import (
	"context"
	"crypto/ecdsa"

	"github.com/dogefuzz/dogefuzz/pkg/common"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

type Agent interface {
	Send(ctx context.Context, nonce uint64, contract *common.Contract, functionName string, args ...interface{}) (string, error)
	GetNonce(ctx context.Context) (uint64, error)
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
