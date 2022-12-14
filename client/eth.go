package client

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/taorzhang/toolkit/client/jsonrpc"
	"github.com/taorzhang/toolkit/types/block"
	"math/big"
	"strings"
)

type EthClientType string

const (
	GethType   = EthClientType("geth")
	ErigonType = EthClientType("erigon")
)

type Eth struct {
	client           *jsonrpc.Client
	internalTxClient *jsonrpc.Client
}

func NewEthClient(client *jsonrpc.Client, internalTxClient *jsonrpc.Client) Provider {
	return &Eth{client: client, internalTxClient: internalTxClient}
}

// BlockNumber 获取最新区块高度
func (e *Eth) BlockNumber(ctx context.Context) (uint64, error) {
	var number math.HexOrDecimal64
	err := e.client.Call("eth_blockNumber", &number)
	return uint64(number), err
}

// MethodCall 执行一个eth call
func (e *Eth) MethodCall(ctx context.Context, out interface{}, args ...interface{}) error {
	return e.client.Call("eth_call", &out, args...)
}

// BlockByHash 根据区块hash获取整个区块信息
func (e *Eth) BlockByHash(ctx context.Context, hash block.Hash, full bool) (*block.Block, error) {
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

// GetNonce 获取账户的交易nonce
func (e *Eth) GetNonce(ctx context.Context, addr block.Address, status string) (nonce uint64, err error) {
	var result string
	err = e.client.Call("eth_getTransactionCount", &result, addr.String(), status)
	if err != nil {
		return
	}
	if !strings.HasPrefix(strings.ToLower(result), "0x") {
		return 0, fmt.Errorf("result '%s' is not correct", result)
	}
	var bigInt big.Int
	data, ok := bigInt.SetString(result[2:], 16)
	if !ok {
		return 0, fmt.Errorf("result '%s' is not correct", result)
	}
	return data.Uint64(), nil
}

// SendTx 发送一笔交易
func (e *Eth) SendTx(ctx context.Context, signTx string) (result string, err error) {
	err = e.client.Call("eth_sendRawTransaction", &result, signTx)
	return
}

// BalanceAt 查询eth余额
func (e *Eth) BalanceAt(ctx context.Context, address block.Address) (*big.Int, error) {
	var result string
	err := e.client.Call("eth_getBalance", &result, address.String(), "latest")
	if err != nil {
		return nil, err
	}
	var i block.BigInt
	err = i.Scan(result)
	if err != nil {
		return nil, err
	}
	return i.ToBigInt(), nil
}

// GetGasPrice 获取gas price
func (e *Eth) GetGasPrice(ctx context.Context) (*big.Int, error) {
	var result string
	err := e.client.Call("eth_gasPrice", &result)
	if err != nil {
		return nil, err
	}
	if !strings.HasPrefix(strings.ToLower(result), "0x") {
		return nil, fmt.Errorf("result '%s' is not correct", result)
	}
	var bigInt big.Int
	data, ok := bigInt.SetString(result[2:], 16)
	if !ok {
		return nil, fmt.Errorf("result '%s' is not correct", result)
	}
	return data, nil
}

// BlockByNumber 根据区块高度获取区块信息
func (e *Eth) BlockByNumber(ctx context.Context, height uint64, full bool) (*block.Block, error) {
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
func (e *Eth) BlocksByNumbers(ctx context.Context, heights []uint64, full bool) ([]*block.Block, error) {
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
				blocks[blockIdx].Transactions[txIdx].InternalTraceCalls = make([]*block.InternalTxCallTrace, 0)
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
func (e *Eth) TransactionByHash(ctx context.Context, hash block.Hash, full bool) (*block.Transaction, error) {
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
func (e *Eth) TransactionsByHashList(ctx context.Context, hashList []block.Hash, full bool) ([]*block.Transaction, error) {
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

// EstimateGas 矿工费估算
func (e *Eth) EstimateGas(ctx context.Context, call CallParameter) (*big.Int, error) {
	var result string
	err := e.client.Call("eth_estimateGas", &result, call.ToArg())
	if err != nil {
		return nil, err
	}
	return block.HexStrToBigInt(result)
}

// GasTipCap 矿工费预估
func (e *Eth) GasTipCap(ctx context.Context) (*big.Int, error) {
	var result string
	err := e.client.Call("eth_maxPriorityFeePerGas", &result)
	if err != nil {
		return nil, err
	}
	return block.HexStrToBigInt(result)
}

// ChainID 获取链ID
func (e *Eth) ChainID(ctx context.Context) (*big.Int, error) {
	var chainID math.HexOrDecimal64
	err := e.client.Call("eth_chainId", &chainID)
	if err != nil {
		return nil, err
	}
	return big.NewInt(int64(chainID)), nil
}

func (e *Eth) InternalTxs(ctx context.Context, txHashes []block.Hash, clientType EthClientType) (map[string][]*block.InternalTxCallTrace, error) {
	switch clientType {
	case ErigonType:
		return e.erigonInternalTx(ctx, txHashes)
	case GethType:
		return e.gethInternalTx(ctx, txHashes)
	}
	return nil, fmt.Errorf("clientType [%s] is not support", clientType)
}

type ErigonCallerTrace struct {
	Action struct {
		From     block.Hex `json:"from"`
		CallType string    `json:"callType"`
		To       block.Hex `json:"to"`
		Value    string    `json:"value"`
	} `json:"action"`
	TransactionHash block.Hex `json:"transactionHash"`
	BlockNumber     uint64    `json:"blockNumber"`
	Result          struct {
		Address block.Hex `json:"address"`
	}
	Type string `json:"type"`
}

func (e *Eth) erigonInternalTx(ctx context.Context, txHashes []block.Hash) (map[string][]*block.InternalTxCallTrace, error) {
	var result = make(map[string][]*block.InternalTxCallTrace)
	elems := make([]rpc.BatchElem, 0)
	internalTxs := make([][]ErigonCallerTrace, len(txHashes))
	for idx := range txHashes {
		internalTxs[idx] = make([]ErigonCallerTrace, 0)
		elems = append(elems, rpc.BatchElem{
			Method: "trace_transaction",
			Args:   []interface{}{txHashes[idx].String()},
			Result: &internalTxs[idx],
		})
	}
	if err := e.internalTxClient.BatchCall(elems, true); err != nil {
		return nil, err
	}
	for _, internalTx := range internalTxs {
		if len(internalTx) == 0 {
			continue
		}
		for _, callTrance := range internalTx {
			if _, ok := result[callTrance.TransactionHash.Hex()]; !ok {
				result[callTrance.TransactionHash.Hex()] = make([]*block.InternalTxCallTrace, 0)
			}
			if strings.EqualFold(callTrance.Type, "create") {
				result[callTrance.TransactionHash.Hex()] = append(result[callTrance.TransactionHash.Hex()], &block.InternalTxCallTrace{
					BlockNumber:     callTrance.BlockNumber,
					TxHash:          callTrance.TransactionHash,
					From:            callTrance.Action.From,
					ContractAddress: callTrance.Result.Address,
				})
			} else {
				if len(callTrance.Action.Value) > 2 && !strings.EqualFold(callTrance.Action.Value, "0x0") {
					bigInt, _ := block.HexStrToBigInt(callTrance.Action.Value)
					result[callTrance.TransactionHash.Hex()] = append(result[callTrance.TransactionHash.Hex()], &block.InternalTxCallTrace{
						BlockNumber: callTrance.BlockNumber,
						TxHash:      callTrance.TransactionHash,
						From:        callTrance.Action.From,
						To:          callTrance.Action.To,
						Value:       block.BigInt(*bigInt),
					})
				}
			}
		}
	}
	return result, nil
}

func (e *Eth) gethInternalTx(ctx context.Context, txHashes []block.Hash) (map[string][]*block.InternalTxCallTrace, error) {
	var result = make(map[string][]*block.InternalTxCallTrace)

	return result, nil
}
