package abi

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/taorzhang/toolkit/abi"
	"github.com/taorzhang/toolkit/client"
	"github.com/taorzhang/toolkit/client/jsonrpc"
	"github.com/taorzhang/toolkit/types/block"
	"math/big"
	"testing"
)

func initProvider(t *testing.T) client.Provider {
	opts := jsonrpc.GetDefaultOpts("https://arbitrum-mainnet.token.im")
	opts = append(opts, jsonrpc.WithRpcHeaders(map[string]string{
		"deviceToken": "test1234",
	}))
	c, err := jsonrpc.NewClient(opts...)
	assert.NoError(t, err)
	return client.NewEthClient(c, nil)
}

func TestABI_Event(t *testing.T) {
	ab, err := abi.NewABIFromList([]string{
		"event TransferBatch(address indexed operator, address indexed from, address indexed to, uint256[] ids, uint256[] values)",
	})
	assert.NoError(t, err)
	event := ab.GetEventByID("0x4a39dc06d4c0dbc64b70af90fd698a233a518aa5d07e595d983b8c0526c8f7fb")
	assert.NotNil(t, event)
	t.Log("id", event.ID())

	provider := initProvider(t)
	transaction, err := provider.TransactionByHash(context.Background(), block.Hex2Hash("0x00000C558A0C8D80A2DE98B3D634FDAA6ECF9DF686B4157F090DA28D4D19E85E"), true)
	assert.NoError(t, err)
	log := transaction.Receipt.Logs[1]
	t.Log("idx", transaction.Hash)
	t.Log("log", log)

	parsedData, err := event.ParseLog(log)
	assert.NoError(t, err)
	t.Log("parsed result", parsedData)

}

func TestABI_2(t *testing.T) {
	newInt := big.NewInt(64)
	t.Log("len", newInt.BitLen())
}

func TestABI_3(t *testing.T) {
	method, err := abi.NewMethod("function transfer(address dst, uint256 wad)")
	assert.NoError(t, err)
	t.Log("method", string(method.HexID()))
}

func TestABI_encode1(t *testing.T) {
	ab, err := abi.NewABIFromList([]string{
		"function echo(address from, address to, string words) public returns(string)",
	})
	assert.NoError(t, err)
	method := ab.GetMethodByID("0xc62969e0")
	encode, err := method.Encode([]interface{}{"0x4a2328fd4790f0950ddaf2f8a369786f094a1299", "0x68b3465833fb72a70ecdf485e0e4c7bd8665fc45", "hello world! Nick!"})
	assert.NoError(t, err)
	var h block.Hex
	err = h.ToHex(encode)
	assert.NoError(t, err)
	t.Log("encode:", h.String())
}

// TestABI_transfer_method 构造input
// https://goerli.etherscan.io/tx/0x669cd6d0afc326785e668ad207b43007be678d05bab5d6d13b240b1379af79f0
func TestABI_transfer_method(t *testing.T) {
	method, err := abi.NewMethod("function transfer(address dst, uint256 wad)")
	assert.NoError(t, err)
	t.Log("sig:", method.Sig())
	t.Log("id:", block.Hex(method.ID()).Hex())
	encode, err := method.Encode([]interface{}{"0x7a64BE40B9f2412FBaDeB519B2d119eC59D8e0CD", "1000000000000000000000"})
	assert.NoError(t, err)
	t.Log("encode", block.Hex(encode).Hex())
	assert.Equal(t,
		"0xa9059cbb0000000000000000000000007a64be40b9f2412fbadeb519b2d119ec59d8e0cd00000000000000000000000000000000000000000000003635c9adc5dea00000",
		block.Hex(encode).Hex())
}
