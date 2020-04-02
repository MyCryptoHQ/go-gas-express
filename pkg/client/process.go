package client

import (
	ClientTypes "github.com/ethereum/go-ethereum/core/types"
)

func FetchBlockMinGas(txs ClientTypes.Transactions) int64 {
	var minGasPrice int64

	for _, tx := range txs {
		if tx.GasPrice().Int64() < minGasPrice || minGasPrice == 0 {
			minGasPrice = tx.GasPrice().Int64()
		}
	}
	if minGasPrice == 0 {
		minGasPrice = 100000000
	}

	return minGasPrice
}
