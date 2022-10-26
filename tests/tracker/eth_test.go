package tracker

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/taorzhang/toolkit/client"
	jsonrpc2 "github.com/taorzhang/toolkit/client/jsonrpc"
	"github.com/taorzhang/toolkit/types/block"
	"testing"
	"time"
)

// initArbTracker 初始化arb tracker
func initArbTracker(t *testing.T) client.Provider {
	opts := jsonrpc2.GetEthCfgOpts(
		"https://arbitrum-mainnet.token.im", 5, 100, 20, 5*time.Second)
	c, err := jsonrpc2.NewClient(
		map[string]string{},
		opts...)
	assert.NoError(t, err)
	return client.NewEthClient(c, nil)
}

// initEthTracker 初始化eth tracker
func initEthTracker(t *testing.T) client.Provider {
	opts := jsonrpc2.GetEthCfgOpts(
		"https://mainnet-eth.token.im", 5, 100, 20, 5*time.Second)
	c, err := jsonrpc2.NewClient(
		map[string]string{},
		opts...)
	assert.NoError(t, err)
	return client.NewEthClient(c, nil)
}

func initEthInternalTracker(t *testing.T) client.Provider {
	opts := jsonrpc2.GetEthCfgOpts("http://goerli.testnet.private:8080/openethereum", 5, 20, 5, 5*time.Second)
	c, err := jsonrpc2.NewClient(
		map[string]string{},
		opts...)
	assert.NoError(t, err)
	assert.NoError(t, err)
	return client.NewEthClient(c, c)
}

func Test_Eth_GetNumber(t *testing.T) {
	provider := initArbTracker(t)
	number, err := provider.BlockNumber(context.Background())
	assert.NoError(t, err)
	t.Log("number:", number)
}

func Test_Eth_ChainID(t *testing.T) {
	provider := initArbTracker(t)
	chainID, err := provider.ChainID(context.Background())
	assert.NoError(t, err)
	t.Log("chainID", chainID)
}

func Test_Eth_GetBlock(t *testing.T) {
	provider := initEthTracker(t)
	number, err := provider.BlockNumber(context.Background())
	assert.NoError(t, err)
	blockData, err := provider.BlockByNumber(context.Background(), number, true)
	assert.NoError(t, err)
	for _, tx := range blockData.Transactions {
		t.Log("idx", tx.TransactionIndex, "hash", tx.Hash)
	}
}

func Test_Eth_GetBlocks(t *testing.T) {
	provider := initEthTracker(t)
	blocks, err := provider.BlocksByNumbers(context.Background(), []uint64{15422138, 15422137}, true)
	assert.NoError(t, err)
	for _, blockData := range blocks {
		t.Log("block number", blockData.Number)
		for _, tx := range blockData.Transactions {
			t.Log("tx_idx", tx.TransactionIndex, "tx_hash", tx.Hash)
		}
	}
}

func Test_Eth_GetBlockByHash(t *testing.T) {
	provider := initEthTracker(t)
	blockData, err := provider.BlockByHash(context.Background(), block.Hex2Hash("0x5d257d69c5b97be15e4aad5ef86e173ea6fc9d16adc5f0fb71f8648e7522c52d"), true)
	assert.NoError(t, err)
	t.Log("tx count", len(blockData.Transactions))
	for _, tx := range blockData.Transactions {
		t.Log("tx_idx", tx.TransactionIndex, "tx_hash", tx.Hash)
	}
}

func Test_eth_internal_Txs(t *testing.T) {
	provider := initEthInternalTracker(t)
	txHashes := []block.Hash{
		block.Hex2Hash("0x48366dce1034414f5fe7b180140d8a78d9d7d83b10774b2901c705191050fd4f"),
		block.Hex2Hash("0x8d8a524da2e9cb686b9d5f7c47af74908238303cea42149b2e316ccc72743ac4"),
	}
	txs, err := provider.InternalTxs(context.Background(), txHashes, client.ErigonType)
	assert.NoError(t, err)
	t.Log("txs", txs)
	for txHash := range txs {
		fmt.Println("tx_hash", txHash)
		for _, internalTx := range txs[txHash] {
			t.Log("From", internalTx.From, "to", internalTx.To, "contract_address", internalTx.ContractAddress, "value", internalTx.Value, "blockNumber", internalTx.BlockNumber)
		}
	}

}
