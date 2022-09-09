package client

import (
	"context"
	"github.com/taorzhang/toolkit/types/block"
	"math/big"
)

type Provider interface {
	ChainID(ctx context.Context) (*big.Int, error)
	BlockNumber(ctx context.Context) (uint64, error)
	GetGasPrice(ctx context.Context) (*big.Int, error)
	SendTx(ctx context.Context, signTx string) (result string, err error)
	EstimateGas(ctx context.Context, call CallParameter) (*big.Int, error)
	BalanceAt(ctx context.Context, address block.Address) (*big.Int, error)
	BlockByHash(ctx context.Context, hash block.Hash, full bool) (*block.Block, error)
	BlockByNumber(ctx context.Context, height uint64, full bool) (*block.Block, error)
	BlocksByNumbers(ctx context.Context, heights []uint64, full bool) ([]*block.Block, error)
	GetNonce(ctx context.Context, addr block.Address, status string) (nonce uint64, err error)
	TransactionByHash(ctx context.Context, hash block.Hash, full bool) (*block.Transaction, error)
	TransactionsByHashList(ctx context.Context, hash []block.Hash, full bool) ([]*block.Transaction, error)
}

var defaultGasLimit = "0x30000"

type CallParameter struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Data     string `json:"data"`
	Gas      string `json:"gas"`
	GasPrice string `json:"gasPrice"`
	Value    string `json:"value"`
}

func (c CallParameter) ToArg() interface{} {
	arg := make(map[string]interface{})
	if c.From != "" {
		arg["from"] = c.From
	}
	if c.To != "" {
		arg["to"] = c.To
	}
	if c.Data != "" {
		arg["data"] = c.Data
	}
	if c.Value != "" {
		arg["value"] = c.Value
	}
	if c.Gas != "" {
		arg["gas"] = c.Gas
	}
	if c.GasPrice != "" {
		arg["gasPrice"] = c.GasPrice
	}
	return arg
}
