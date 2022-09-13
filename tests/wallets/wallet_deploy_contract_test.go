package wallets

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var erc721ContractAddress = "0x924C9C38Bac8e124fc0D2A747a009791Ea436b7c"

// TestNewWallet_deployContract
// 部署合约
// for test: https://goerli.etherscan.io/address/0x924c9c38bac8e124fc0d2a747a009791ea436b7c
// contract address: 0x924C9C38Bac8e124fc0D2A747a009791Ea436b7c
func TestNewWallet_deployContract(t *testing.T) {
	wallet := initWallet(t)
	content, err := os.ReadFile("./contracts/721/output/MqyFt.bin")
	assert.NoError(t, err)
	code := string(content)
	txHash, err := wallet.DeployContract(code)
	if err != nil {
		t.Log(err)
		return
	}
	assert.NoError(t, err)
	t.Log("tx_hash:", txHash.String())
}
