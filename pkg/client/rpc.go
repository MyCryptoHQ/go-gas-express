package client

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	ClientTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
)

func GetBalance(client ethclient.Client, address string) (*big.Int, error) {
	account := common.HexToAddress(address)
	balance, err := client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		return big.NewInt(0), err
	}
	return balance, nil
}

func GetBlockByNum(client ethclient.Client, blockNum int64) (ClientTypes.Block, error) {
	block, err := client.BlockByNumber(context.Background(), big.NewInt(blockNum))
	if err != nil {
		return ClientTypes.Block{}, err
	}
	return *block, nil
}

func GetBlockNumber(client ethclient.Client) (int64, error) {
	blockNumHeader, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return int64(0), err
	}
	return blockNumHeader.Number.Int64(), nil
}
