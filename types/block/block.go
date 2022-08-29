package block

import (
	"github.com/ethereum/go-ethereum/common/math"
)

type Block struct {
	Number             math.HexOrDecimal64
	Hash               Hash
	ParentHash         Hash
	Sha3Uncles         Hash
	TransactionsRoot   Hash
	StateRoot          Hash
	ReceiptsRoot       Hash
	Miner              Address
	Difficulty         *BigInt
	ExtraData          Hex
	GasLimit           math.HexOrDecimal64
	GasUsed            math.HexOrDecimal64
	Timestamp          math.HexOrDecimal64
	Transactions       []*Transaction
	TransactionsHashes []Hash
	Uncles             []Hash
}
