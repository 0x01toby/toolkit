package polling

import (
	"context"
	"github.com/taorzhang/toolkit/client"
	"github.com/taorzhang/toolkit/types/block"
)

type Item struct {
	ctx                context.Context
	blocks             []*block.Block
	client             client.Provider
	start              uint
	end                uint
	skipInternal       bool
	internalClientType client.EthClientType
}

func NewItem(ctx context.Context, client client.Provider, start, end uint, skipInternal bool, internalClientType client.EthClientType) *Item {
	return &Item{ctx: ctx, blocks: make([]*block.Block, 0), start: start, end: end, client: client, skipInternal: skipInternal, internalClientType: internalClientType}
}

func (i *Item) Retrieve() error {
	var heights []uint64
	for height := i.start; height < i.end; height++ {
		heights = append(heights, uint64(height))
	}
	blocks, err := i.client.BlocksByNumbers(i.ctx, heights, true)
	if err != nil {
		return err
	}
	if i.skipInternal {
		i.blocks = append(i.blocks, blocks...)
		return nil
	}
	var txHashes []block.Hash
	for idx := range blocks {
		txHashes = append(txHashes, blocks[idx].Hash)
	}
	internalTxs, err := i.client.InternalTxs(i.ctx, txHashes, i.internalClientType)
	if err != nil {
		return err
	}
	for blockIdx := range blocks {
		for txIdx := range blocks[blockIdx].Transactions {
			if internalCallTraces, ok := internalTxs[blocks[blockIdx].Transactions[txIdx].Hash.String()]; ok {
				blocks[blockIdx].Transactions[txIdx].InternalTraceCalls = append(blocks[blockIdx].Transactions[txIdx].InternalTraceCalls, internalCallTraces...)
			}
		}
	}
	i.blocks = append(i.blocks, blocks...)
	return nil
}
