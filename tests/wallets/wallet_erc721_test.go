package wallets

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/taorzhang/toolkit/abi"
	"github.com/taorzhang/toolkit/client"
	"github.com/taorzhang/toolkit/types/block"
	wallet2 "github.com/taorzhang/toolkit/wallet"
	"math/big"
	"testing"
)

// TestNewWallet_erc721_symbol
// 查询symbol
func TestNewWallet_erc721_symbol(t *testing.T) {
	wallet := initWallet(t)
	method, err := abi.NewMethod("function symbol()")
	assert.NoError(t, err)
	var symbol string
	err = wallet.Client.MethodCall(context.Background(), &symbol, client.CallParameter{
		From: *block.EmptyAddress.ToCommonAddress(),
		To:   block.Hexstr2Address(erc721ContractAddress).ToCommonAddress(),
		Data: common.FromHex(block.Hex(method.ID()).Hex()),
	}.ToArg(), "latest")
	assert.NoError(t, err)
	t.Log("symbol:", block.HexstrToString(symbol))
	assert.Equal(t, "MQY", block.HexstrToString(symbol))
}

// TestNewWallet_erc721_name
// 查询name
func TestNewWallet_erc721_name(t *testing.T) {
	wallet := initWallet(t)
	method, err := abi.NewMethod("function name()")
	assert.NoError(t, err)
	var symbol string
	err = wallet.Client.MethodCall(context.Background(), &symbol, client.CallParameter{
		From: *block.EmptyAddress.ToCommonAddress(),
		To:   block.Hexstr2Address(erc721ContractAddress).ToCommonAddress(),
		Data: common.FromHex(block.Hex(method.ID()).Hex()),
	}.ToArg(), "latest")
	assert.NoError(t, err)
	t.Log("name:", block.HexstrToString(symbol))
	assert.Equal(t, "MQY NFT", block.HexstrToString(symbol))
}

// TestNewWallet_erc721_mint
// mint 721 token
func TestNewWallet_erc721_mint(t *testing.T) {
	wallet := initWallet(t)
	method, err := abi.NewMethod("function mint(address _to, uint256 _tokenId, string _uri)")
	assert.NoError(t, err)
	contract := block.Hexstr2Address(erc721ContractAddress)
	toAddress := block.Hexstr2Address("0x9236B49DA606d83b3c69004D13fd14f9F545A90B")
	txData, err := wallet.CreateLegacyTxData(wallet2.Pending, contract, big.NewInt(0), method, toAddress.String(), big.NewInt(1), "https://www.baidu.com/{1}")
	assert.NoError(t, err)
	signedTx, err := wallet.SignTx(txData)
	assert.NoError(t, err)
	tx, err := wallet.Client.SendTx(context.Background(), signedTx)
	assert.NoError(t, err)
	fmt.Println("tx_hash:", tx)
	// like this: https://goerli.etherscan.io/tx/0x643706e02db973ae88273e98bbbf6e4c4dda44e1ef7b94574256160f753804a5
}

// TestNewWallet_erc721_transfer
// 发送 721 token
func TestNewWallet_erc721_transfer(t *testing.T) {
	wallet := initWallet2(t)
	contract := block.Hexstr2Address(erc721ContractAddress)
	toAddress := block.Hexstr2Address("0x55D65F2dE30632e224766CF6652E02d5753B0fda")
	txHash, err := wallet.SendErc721(contract, toAddress, big.NewInt(1))
	assert.NoError(t, err)
	t.Log("tx_hash:", txHash)
	// like this: https://goerli.etherscan.io/tx/0xe91a20a1a463b8a54e663b45609363e340dff62093fc43ea01507d03c45c418f
}

// TestNewWallet_erc721_ownerOf
// 查询tokenID的owner
func TestNewWallet_erc721_ownerOf(t *testing.T) {
	wallet := initWallet(t)
	ownerOf, err := wallet.Erc721OwnerOf(block.Hexstr2Address(erc721ContractAddress), big.NewInt(1))
	assert.NoError(t, err)
	t.Log("owner:", ownerOf)
	assert.Equal(t, "0x55D65F2dE30632e224766CF6652E02d5753B0fda", ownerOf.String())
}
