package interfaces

import (
	"context"
	"crypto/ecdsa"

	"github.com/dogefuzz/dogefuzz/pkg/common"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

type Agent interface {
	Send(ctx context.Context, contract *common.Contract, functionName string, args ...interface{}) (string, error)
}

type Deployer interface {
	Deploy(ctx context.Context, contract *common.Contract, args ...interface{}) (string, error)
}

type Wallet interface {
	GetPrivateKey() *ecdsa.PrivateKey
	GetPrivateKeyHex() string
	GetPublicKey() *ecdsa.PublicKey
	GetPublicKeyHex() string
	GetAddress() gethcommon.Address
}
