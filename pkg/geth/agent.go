package geth

import (
	"context"
	"math/big"
	"strings"

	"github.com/dogefuzz/dogefuzz/config"
	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/dogefuzz/dogefuzz/pkg/interfaces"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type agent struct {
	client *ethclient.Client
	wallet interfaces.Wallet
	cfg    config.GethConfig
}

func NewAgent(cfg config.GethConfig) (*agent, error) {
	wallet, err := NewWalletFromPrivateKeyHex(cfg.AgentPrivateKeyHex)
	if err != nil {
		return nil, err
	}

	client, err := ethclient.Dial(cfg.NodeAddress)
	if err != nil {
		return nil, err
	}

	return &agent{client, wallet, cfg}, nil
}

func (d *agent) Send(ctx context.Context, nonce uint64, contract *common.Contract, functionName string, args ...interface{}) (string, error) {
	_ = nonce

	parsedABI, err := abi.JSON(strings.NewReader(contract.AbiDefinition))
	if err != nil {
		return "", err
	}

	boundContract := bind.NewBoundContract(gethcommon.HexToAddress(contract.Address), parsedABI, d.client, d.client, d.client)

	gasPrice, err := d.client.SuggestGasPrice(ctx)
	if err != nil {
		return "", err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(d.wallet.GetPrivateKey(), big.NewInt(d.cfg.ChainID))
	if err != nil {
		return "", err
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(2000000)
	auth.GasPrice = gasPrice

	tx, err := boundContract.Transact(auth, functionName, args...)
	if err != nil {
		return "", err
	}

	return tx.Hash().Hex(), nil
}

func (a *agent) GetNonce(ctx context.Context) (uint64, error) {
	nonce, err := a.client.PendingNonceAt(ctx, a.wallet.GetAddress())
	if err != nil {
		return 0, err
	}
	return nonce, err
}
