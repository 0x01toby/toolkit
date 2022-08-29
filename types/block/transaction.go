package block

import (
	"github.com/ethereum/go-ethereum/common/math"
)

type Transaction struct {
	// see const
	Type math.HexOrDecimal64

	// legacy values
	Hash     Hash
	From     Address
	To       *Address
	Input    Hex
	GasPrice math.HexOrDecimal64
	Gas      math.HexOrDecimal64
	Value    *BigInt
	Nonce    math.HexOrDecimal64
	V        math.HexOrDecimal256
	R        Hex
	S        Hex

	// jsonrpc values
	BlockHash        Hash
	BlockNumber      math.HexOrDecimal64
	TransactionIndex math.HexOrDecimal64

	// eip-2930 values
	ChainID    *math.HexOrDecimal64
	AccessList AccessList

	// eip-1559 values
	MaxPriorityFeePerGas *BigInt
	MaxFeePerGas         *BigInt

	// receipt
	Receipt *Receipt
}
