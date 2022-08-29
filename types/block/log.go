package block

import "github.com/ethereum/go-ethereum/common/math"

type Log struct {
	Removed          bool
	LogIndex         math.HexOrDecimal64
	TransactionIndex math.HexOrDecimal64
	TransactionHash  Hash
	BlockHash        Hash
	BlockNumber      math.HexOrDecimal64
	Address          Address
	Topics           []Hash
	Data             Hex
}
