package wallets

import (
	"github.com/stretchr/testify/assert"
	"github.com/taorzhang/toolkit/wallet"
	"testing"
)

func TestHdWallet_case1(t *testing.T) {
	hd, err := wallet.NewHDWalletFromMnemonic("tag volcano eight thank tide danger coast health above argue embrace heavy")
	assert.NoError(t, err)
	pk, err := hd.GeneratePrivateKeyByPath("m/44'/60'/0'/0/0")
	assert.NoError(t, err)
	assert.Equal(t, "0xC49926C4124cEe1cbA0Ea94Ea31a6c12318df947", pk.PublicAddress())
	subPk, err := hd.GeneratePrivateKeyByPath("m/44'/60'/0'/0/1")
	assert.NoError(t, err)
	assert.Equal(t, "0x8230645aC28A4EdD1b0B53E7Cd8019744E9dD559", subPk.PublicAddress())
	assert.Equal(t, "tag volcano eight thank tide danger coast health above argue embrace heavy", hd.Mnemonic())
}
