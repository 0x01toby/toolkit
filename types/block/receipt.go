package block

import "github.com/ethereum/go-ethereum/common/math"

type Receipt struct {
	TransactionHash   Hash
	TransactionIndex  math.HexOrDecimal64
	ContractAddress   Address
	BlockHash         Hash
	From              Address
	BlockNumber       math.HexOrDecimal64
	GasUsed           math.HexOrDecimal64
	CumulativeGasUsed math.HexOrDecimal64
	LogsBloom         Hex
	Logs              []*Log
	Status            math.HexOrDecimal64
}
