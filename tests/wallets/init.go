package wallets

import (
	"github.com/stretchr/testify/assert"
	"github.com/taorzhang/toolkit/client"
	jsonrpc2 "github.com/taorzhang/toolkit/client/jsonrpc"
	wallet2 "github.com/taorzhang/toolkit/wallet"
	"os"
	"testing"
)

func initProvider(t *testing.T) client.Provider {
	opts := jsonrpc2.GetDefaultOpts("https://goerli.infura.io/v3/fa57784b5b834db1b685341ec9867a3a")
	c, err := jsonrpc2.NewClient(opts...)
	assert.NoError(t, err)
	return client.NewEthClient(c, nil)
}

func initWallet(t *testing.T) *wallet2.Account {
	content, err := os.ReadFile("./privatekey.env")
	assert.NoError(t, err)
	wallet := wallet2.NewWalletFromPrivateKeyStr(string(content), wallet2.WithEthProvider(initProvider(t)))
	return wallet
}

func initWallet2(t *testing.T) *wallet2.Account {
	content, err := os.ReadFile("./keystore.json")
	assert.NoError(t, err)
	password, err := os.ReadFile("./keystore.pass.env")
	assert.NoError(t, err)
	store := wallet2.NewWalletFromKeyStore(content, string(password), wallet2.WithEthProvider(initProvider(t)))
	return store
}
