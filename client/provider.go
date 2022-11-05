package client

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/taorzhang/toolkit/types/block"
	"math/big"
)

type Provider interface {
	ChainID(ctx context.Context) (*big.Int, error)
	BlockNumber(ctx context.Context) (uint64, error)
	GasTipCap(ctx context.Context) (*big.Int, error)
	GetGasPrice(ctx context.Context) (*big.Int, error)
	SendTx(ctx context.Context, signTx string) (result string, err error)
	EstimateGas(ctx context.Context, call CallParameter) (*big.Int, error)
	BalanceAt(ctx context.Context, address block.Address) (*big.Int, error)
	MethodCall(ctx context.Context, out interface{}, args ...interface{}) error
	BlockByHash(ctx context.Context, hash block.Hash, full bool) (*block.Block, error)
	BlockByNumber(ctx context.Context, height uint64, full bool) (*block.Block, error)
	BlocksByNumbers(ctx context.Context, heights []uint64, full bool) ([]*block.Block, error)
	GetNonce(ctx context.Context, addr block.Address, status string) (nonce uint64, err error)
	TransactionByHash(ctx context.Context, hash block.Hash, full bool) (*block.Transaction, error)
	TransactionsByHashList(ctx context.Context, hash []block.Hash, full bool) ([]*block.Transaction, error)
	InternalTxs(ctx context.Context, txHashes []block.Hash, clientType EthClientType) (map[string][]*block.InternalTxCallTrace, error)
}

var DefaultGasLimit = "0x30000"
var DefaultGasLimitInt uint64 = 30000

type CallParameter struct {
	From     common.Address
	To       *common.Address
	Data     []byte
	Gas      uint64
	GasPrice *big.Int
	Value    *big.Int
}

func (c CallParameter) ToArg() interface{} {
	arg := make(map[string]interface{})
	arg["from"] = c.From
	if c.To != nil {
		arg["to"] = c.To
	}
	if c.Data != nil {
		arg["data"] = hexutil.Bytes(c.Data)
	}
	if c.Value != nil {
		arg["value"] = (*hexutil.Big)(c.Value)
	}
	if c.Gas != 0 {
		arg["gas"] = hexutil.Uint64(c.Gas)
	}
	if c.GasPrice != nil {
		arg["gasPrice"] = (*hexutil.Big)(c.GasPrice)
	}
	return arg
}
