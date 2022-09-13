package tracker

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/taorzhang/toolkit/client"
	"github.com/taorzhang/toolkit/jsonrpc"
	"github.com/taorzhang/toolkit/types/block"
	"testing"
	"time"
)

// initArbTracker 初始化arb tracker
func initArbTracker(t *testing.T) client.Provider {
	opts := jsonrpc.GetEthCfgOpts(
		"https://arbitrum-mainnet.token.im", 5, 100, 20, 5*time.Second)
	c, err := jsonrpc.NewClient(
		map[string]string{},
		opts...)
	assert.NoError(t, err)
	return client.NewEthClient(c)
}

// initEthTracker 初始化eth tracker
func initEthTracker(t *testing.T) client.Provider {
	opts := jsonrpc.GetEthCfgOpts(
		"https://mainnet-eth.token.im", 5, 100, 20, 5*time.Second)
	c, err := jsonrpc.NewClient(
		map[string]string{},
		opts...)
	assert.NoError(t, err)
	return client.NewEthClient(c)
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
