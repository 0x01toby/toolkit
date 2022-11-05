package block

type InternalTxCallTrace struct {
	From            Hex
	To              Hex
	ContractAddress Hex
	Value           BigInt
	TxHash          Hex
	BlockNumber     uint64
}
