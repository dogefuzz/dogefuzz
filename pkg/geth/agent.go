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
	cfg    config.GethConfig
}

func NewAgent(cfg config.GethConfig) (*agent, error) {

	client, err := ethclient.Dial(cfg.NodeAddress)
	if err != nil {
		return nil, err
	}

	return &agent{client, cfg}, nil
}

func (d *agent) Send(ctx context.Context, wallet interfaces.Wallet, contract *common.Contract, functionName string, value *big.Int, args ...interface{}) (string, error) {
	parsedABI, err := abi.JSON(strings.NewReader(contract.AbiDefinition))
	if err != nil {
		return "", err
	}
	addressAsHexString := gethcommon.HexToAddress(contract.Address)
	boundContract := bind.NewBoundContract(addressAsHexString, parsedABI, d.client, d.client, d.client)

	nonce, err := d.client.PendingNonceAt(ctx, wallet.GetAddress())
	if err != nil {
		return "", err
	}

	gasPrice, err := d.client.SuggestGasPrice(ctx)
	if err != nil {
		return "", err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(wallet.GetPrivateKey(), big.NewInt(d.cfg.ChainID))
	if err != nil {
		return "", err
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = value
	auth.GasLimit = d.cfg.MinGasLimit
	auth.GasPrice = gasPrice
	auth.Context = ctx

	tx, err := boundContract.Transact(auth, functionName, args...)
	if err != nil {
		return "", err
	}

	return tx.Hash().Hex(), nil
}

func (d *agent) Transfer(ctx context.Context, wallet interfaces.Wallet, contract *common.Contract, value *big.Int) (string, error) {
	parsedABI, err := abi.JSON(strings.NewReader(contract.AbiDefinition))
	if err != nil {
		return "", err
	}
	addressAsHexString := gethcommon.HexToAddress(contract.Address)
	boundContract := bind.NewBoundContract(addressAsHexString, parsedABI, d.client, d.client, d.client)

	nonce, err := d.client.PendingNonceAt(ctx, wallet.GetAddress())
	if err != nil {
		return "", err
	}

	gasPrice, err := d.client.SuggestGasPrice(ctx)
	if err != nil {
		return "", err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(wallet.GetPrivateKey(), big.NewInt(d.cfg.ChainID))
	if err != nil {
		return "", err
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = value
	auth.GasLimit = d.cfg.MinGasLimit
	auth.GasPrice = gasPrice
	auth.Context = ctx

	tx, err := boundContract.Transfer(auth)
	if err != nil {
		return "", err
	}

	return tx.Hash().Hex(), nil
}
