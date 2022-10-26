package block

type InternalTransaction struct {
	From            Hex
	To              Hex
	ContractAddress Hex
	Value           BigInt
	TxHash          Hex
	BlockNumber     uint64
}
