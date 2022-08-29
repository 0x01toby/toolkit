package tracker

import (
	"github.com/taorzhang/toolkit/types/block"
	"math/big"
)

type Provider interface {
	BlockNumber() (uint64, error)
	BlockByHash(hash block.Hash, full bool) (*block.Block, error)
	BlockByNumber(height uint64, full bool) (*block.Block, error)
	BlocksByNumbers(heights []uint64, full bool) ([]*block.Block, error)
	TransactionByHash(hash block.Hash, full bool) (*block.Transaction, error)
	TransactionsByHashList(hash []block.Hash, full bool) ([]*block.Transaction, error)
	ChainID() (*big.Int, error)
}
