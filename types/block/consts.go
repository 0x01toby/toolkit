package block

import "github.com/ethereum/go-ethereum/common/math"

const (
	TransactionLegacy math.HexOrDecimal64 = 0
	// TransactionAccessList eip-2930
	TransactionAccessList math.HexOrDecimal64 = 1
	// TransactionDynamicFee eip-1559
	TransactionDynamicFee math.HexOrDecimal64 = 2
)
