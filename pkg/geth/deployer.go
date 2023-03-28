package geth

import (
	"context"
	"math/big"
	"strings"
	"time"

	"github.com/dogefuzz/dogefuzz/config"
	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/dogefuzz/dogefuzz/pkg/interfaces"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

type deployer struct {
	client *ethclient.Client
	wallet interfaces.Wallet
	cfg    config.GethConfig
	logger *zap.Logger
}

func NewDeployer(logger *zap.Logger, cfg config.GethConfig) (*deployer, error) {
	wallet, err := NewWalletFromPrivateKeyHex(cfg.DeployerPrivateKeyHex)
	if err != nil {
		logger.Sugar().Errorf("failed to import wallet from private key: %v", err)
		return nil, err
	}

	client, err := ethclient.Dial(cfg.NodeAddress)
	if err != nil {
		logger.Sugar().Errorf("failed to connect client with node address %s: %v", cfg.NodeAddress, err)
		return nil, err
	}

	return &deployer{client, wallet, cfg, logger}, nil
}

func (d *deployer) Deploy(ctx context.Context, contract *common.Contract, args ...interface{}) (string, string, error) {
	parsedABI, err := abi.JSON(strings.NewReader(contract.AbiDefinition))
	if err != nil {
		d.logger.Sugar().Errorf("failed to parse contract's ABI definition: %v", err)
		return "", "", err
	}

	nonce, err := d.client.PendingNonceAt(ctx, d.wallet.GetAddress())
	if err != nil {
		d.logger.Sugar().Errorf("failed to get nonce from address %s: %v", d.wallet.GetAddress(), err)
		return "", "", err
	}

	gasPrice, err := d.client.SuggestGasPrice(ctx)
	if err != nil {
		d.logger.Sugar().Errorf("failed to get suggested gas price: %v", err)
		return "", "", err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(d.wallet.GetPrivateKey(), big.NewInt(d.cfg.ChainID))
	if err != nil {
		d.logger.Sugar().Errorf("failed to start request for contract deployment: %v", err)
		return "", "", err
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(2000000)
	auth.GasPrice = gasPrice
	auth.Context = ctx

	_, tx, _, err := bind.DeployContract(auth, parsedABI, gethcommon.FromHex(contract.DeploymentBytecode), d.client, args...)
	if err != nil {
		d.logger.Sugar().Errorf("failed to request contract deployment: %v", err)
		return "", "", err
	}

	var receipt *types.Receipt
	for {
		receipt, err = d.client.TransactionReceipt(ctx, tx.Hash())
		if err != nil {
			if err != ethereum.NotFound {
				d.logger.Sugar().Errorf("failed to get transactio receipt: %v", err)
				return "", "", err
			}
		} else {
			break
		}

		time.Sleep(100 * time.Microsecond)
	}
	contract.Address = receipt.ContractAddress.Hex()
	return receipt.ContractAddress.Hex(), tx.Hash().Hex(), nil
}
