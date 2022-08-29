package tracker

import (
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/taorzhang/toolkit/jsonrpc"
	"github.com/taorzhang/toolkit/types/block"
	"math/big"
)

type Eth struct {
	client *jsonrpc.Client
}

func NewEthTracker(client *jsonrpc.Client) Provider {
	return &Eth{client: client}
}

// BlockNumber 获取最新区块高度
func (e *Eth) BlockNumber() (uint64, error) {
	var number math.HexOrDecimal64
	err := e.client.Call("eth_blockNumber", &number)
	return uint64(number), err
}

// BlockByHash 根据区块hash获取整个区块信息
func (e *Eth) BlockByHash(hash block.Hash, full bool) (*block.Block, error) {
	var blockData block.Block
	err := e.client.Call("eth_getBlockByHash", &blockData, hash, full)
	if err != nil {
		return nil, err
	}
	if full {
		var batch []rpc.BatchElem
		for idx := range blockData.Transactions {
			blockData.Transactions[idx].Receipt = new(block.Receipt)
			batch = append(batch, rpc.BatchElem{
				Method: "eth_getTransactionReceipt",
				Args:   []interface{}{blockData.Transactions[idx].Hash},
				Result: blockData.Transactions[idx].Receipt,
			})
		}
		err = e.client.BatchCall(batch, true)
	}
	return &blockData, err
}

// BlockByNumber 根据区块高度获取区块信息
func (e *Eth) BlockByNumber(height uint64, full bool) (*block.Block, error) {
	var blockData block.Block
	err := e.client.Call("eth_getBlockByNumber", &blockData, math.HexOrDecimal64(height), full)
	if err != nil {
		return nil, err
	}
	if full {
		var batch []rpc.BatchElem
		for idx := range blockData.Transactions {
			blockData.Transactions[idx].Receipt = new(block.Receipt)
			batch = append(batch, rpc.BatchElem{
				Method: "eth_getTransactionReceipt",
				Args:   []interface{}{blockData.Transactions[idx].Hash},
				Result: blockData.Transactions[idx].Receipt,
			})
		}
		err = e.client.BatchCall(batch, true)
	}
	return &blockData, err
}

// BlocksByNumbers 批量获取区块信息
func (e *Eth) BlocksByNumbers(heights []uint64, full bool) ([]*block.Block, error) {
	blocks := make([]*block.Block, len(heights))
	var blockBatch []rpc.BatchElem
	for idx := range heights {
		blocks[idx] = new(block.Block)
		blockBatch = append(blockBatch, rpc.BatchElem{
			Method: "eth_getBlockByNumber",
			Args:   []interface{}{math.HexOrDecimal64(heights[idx]), full},
			Result: blocks[idx],
		})
	}
	err := e.client.BatchCall(blockBatch, true)
	if err != nil {
		return nil, err
	}
	if full {
		var receipts []rpc.BatchElem
		for blockIdx := range blocks {
			for txIdx := range blocks[blockIdx].Transactions {
				blocks[blockIdx].Transactions[txIdx].Receipt = new(block.Receipt)
				receipts = append(receipts, rpc.BatchElem{
					Method: "eth_getTransactionReceipt",
					Args:   []interface{}{blocks[blockIdx].Transactions[txIdx].Hash},
					Result: blocks[blockIdx].Transactions[txIdx].Receipt,
				})
			}
		}
		err = e.client.BatchCall(receipts, true)
	}
	return blocks, err
}

// TransactionByHash 根据交易hash获取交易信息
func (e *Eth) TransactionByHash(hash block.Hash, full bool) (*block.Transaction, error) {
	var tx block.Transaction
	err := e.client.Call("eth_getTransactionByHash", &tx, hash)
	if err != nil {
		return nil, err
	}
	if !full {
		return &tx, nil
	}
	tx.Receipt = new(block.Receipt)
	err = e.client.Call("eth_getTransactionReceipt", tx.Receipt, hash)
	return &tx, err
}

// TransactionsByHashList 批量获取交易信息
func (e *Eth) TransactionsByHashList(hashList []block.Hash, full bool) ([]*block.Transaction, error) {
	txList := make([]*block.Transaction, len(hashList))
	var txBatch []rpc.BatchElem
	for idx := range hashList {
		txList[idx] = new(block.Transaction)
		txList[idx].Receipt = new(block.Receipt)
		txBatch = append(txBatch, rpc.BatchElem{
			Method: "eth_getTransactionByHash",
			Args:   []interface{}{hashList[idx]},
			Result: txList[idx],
		})
	}
	err := e.client.BatchCall(txBatch, true)
	if err != nil {
		return nil, err
	}
	if !full {
		return txList, nil
	}
	var receiptBatch []rpc.BatchElem
	for idx := range txList {
		receiptBatch = append(receiptBatch, rpc.BatchElem{
			Method: "eth_getTransactionReceipt",
			Args:   []interface{}{txList[idx].Hash},
			Result: txList[idx].Receipt,
		})
	}
	err = e.client.BatchCall(receiptBatch, true)
	return txList, err
}

// ChainID 获取链ID
func (e *Eth) ChainID() (*big.Int, error) {
	var chainID math.HexOrDecimal64
	err := e.client.Call("eth_chainId", &chainID)
	if err != nil {
		return nil, err
	}
	return big.NewInt(int64(chainID)), nil
}
