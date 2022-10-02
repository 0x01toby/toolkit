package wallets

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/taorzhang/toolkit/client"
	"github.com/taorzhang/toolkit/types/block"
	wallet2 "github.com/taorzhang/toolkit/wallet"
	"testing"
)

// 0x0e0d81304af3edf2ec843565176316b7456f84f8606a0f0d3c5b46694a617b45
// from 0x55D65F2dE30632e224766CF6652E02d5753B0fda
// to 0x9236B49DA606d83b3c69004D13fd14f9F545A90B
// TestNewWallet_sendNative
// 发送原声代币
func TestNewWallet_sendNative(t *testing.T) {
	wallet := initWallet(t)
	wei, err := wallet2.ToWei("0.1")
	assert.NoError(t, err)
	txHash, err := wallet.SendNativeToken(
		block.Hexstr2Address("0x9236B49DA606d83b3c69004D13fd14f9F545A90B"), wei)
	assert.NoError(t, err)
	t.Log("tx_hash:", txHash)
}

// TestNewWallet_sendErc20
// 发送erc20 token
func TestNewWallet_sendErc20(t *testing.T) {
	wallet := initWallet(t)
	txHash, err := wallet.SendErc20Token(
		block.Hexstr2Address("0x11fE4B6AE13d2a6055C8D9cF65c55bac32B5d844"),
		block.Hexstr2Address("0x9236B49DA606d83b3c69004D13fd14f9F545A90B"),
		"0.5")
	assert.NoError(t, err)
	t.Log("tx_hash:", txHash)
}

// TestNewWallet_approveErc20
// 代币授权
func TestNewWallet_approveErc20(t *testing.T) {
	wallet := initWallet(t)
	txHash, err := wallet.ApprovalErc20Token(
		block.Hexstr2Address("0x11fE4B6AE13d2a6055C8D9cF65c55bac32B5d844"),
		block.Hexstr2Address("0x9236B49DA606d83b3c69004D13fd14f9F545A90B"),
		"100",
	)
	assert.NoError(t, err)
	t.Log("tx_hash:", txHash)
}

// TestNewWallet_transferFromErc20
// approve后，从from Transfer token 到to地址
func TestNewWallet_transferFromErc20(t *testing.T) {
	wallet := initWallet2(t)
	t.Log("address:", wallet.Address())
	txHash, err := wallet.TransferFromErc20Token(
		block.Hexstr2Address("0x11fE4B6AE13d2a6055C8D9cF65c55bac32B5d844"),
		block.Hexstr2Address("0x55D65F2dE30632e224766CF6652E02d5753B0fda"),
		block.EmptyAddress,
		"25.3",
	)
	assert.NoError(t, err)
	t.Log("tx_hash:", txHash)
}

// TestNewWallet_estimateDeployGas
// 预估部署合约的费用
func TestNewWallet_estimateDeployGas(t *testing.T) {
	wallet := initWallet(t)
	code := "6060604052600a8060106000396000f360606040526008565b00"
	nonce, err := wallet.GetNonce("pending")
	assert.NoError(t, err)
	gasPrice, err := wallet.Client.GetGasPrice(context.Background())
	assert.NoError(t, err)
	data, err := wallet2.CreateContractTxData(nonce, client.DefaultGasLimitInt, gasPrice, code)
	assert.NoError(t, err)
	gas, err := wallet.EstimateGas(types.NewTx(data))
	assert.NoError(t, err)
	fmt.Println("gas", gas)
}

// TestNewWallet_sendErc20
// 发送erc20 token
func TestNewWallet_sendErc20_2(t *testing.T) {
	wallet := initWallet2(t)
	txHash, err := wallet.SendErc20Token(
		block.Hexstr2Address("0x11fE4B6AE13d2a6055C8D9cF65c55bac32B5d844"),
		block.Hexstr2Address("0x55D65F2dE30632e224766CF6652E02d5753B0fda"),
		"10")
	assert.NoError(t, err)
	t.Log("tx_hash:", txHash)
}

func TestNewWallet_sendErc20_3(t *testing.T) {
	wallet := initWallet2(t)
	txHash, err := wallet.SendErc20Token(
		block.Hexstr2Address("0xc3359f800Aa2ea472348E26541F55d30E6633243"),
		block.Hexstr2Address("0x6EaebeA65d41354DE3cBdE1DA2d9795ADcd2f875"),
		"99.999999999999999")
	assert.NoError(t, err)
	t.Log("tx_hash:", txHash)
}

func TestNewWallet_mintErc20(t *testing.T) {
	wallet := initWallet2(t)
	txHash, err := wallet.MintErc20Token(
		block.Hexstr2Address("0xc3359f800Aa2ea472348E26541F55d30E6633243"),
		"999")
	assert.NoError(t, err)
	t.Log("tx_hash:", txHash)
}
