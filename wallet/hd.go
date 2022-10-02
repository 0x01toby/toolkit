package wallet

import (
	"crypto/ecdsa"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/tyler-smith/go-bip39"
	"sync"
)

type HDWallet struct {
	masterKey *hdkeychain.ExtendedKey
	entropy   []byte
	paths     map[string]*PrivateKey
	state     sync.RWMutex
}

func NewHDWalletFromMnemonic(mnemonic string) (*HDWallet, error) {
	hdWallet := &HDWallet{paths: map[string]*PrivateKey{}}
	entropy, err := bip39.EntropyFromMnemonic(mnemonic)
	if err != nil {
		return nil, err
	}
	hdWallet.entropy = make([]byte, len(entropy))
	copy(hdWallet.entropy, entropy)
	return hdWallet, nil
}

// Mnemonic 导出助记词
func (hd *HDWallet) Mnemonic() string {
	mnemonic, _ := bip39.NewMnemonic(hd.entropy)
	return mnemonic
}

// GeneratePrivateKeyByPath 根据path生成私钥
func (hd *HDWallet) GeneratePrivateKeyByPath(path string) (*PrivateKey, error) {
	if pk, ok := hd.paths[path]; ok {
		return pk, nil
	}
	hd.state.Lock()
	defer hd.state.Unlock()
	mnemonic, _ := bip39.NewMnemonic(hd.entropy)
	extendKey, err := hd.getExtendKeyByMnemonic(mnemonic, "")
	if err != nil {
		return nil, err
	}
	pk, err := hd.getPkByExtendKeyAndPath(extendKey, path)
	if err != nil {
		return nil, err
	}
	hd.paths[path] = &PrivateKey{PrivateKey: pk}
	return hd.paths[path], nil
}

// getExtendKeyByMnemonic
// 根据助记词获取 extendKey
func (hd *HDWallet) getExtendKeyByMnemonic(mnemonic, password string) (*hdkeychain.ExtendedKey, error) {
	seed, err := bid39Seed(mnemonic, password)
	if err != nil {
		return nil, err
	}
	return hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
}

// getPkByExtendKeyAndPath
// 根据path 和 key生成私钥
func (hd *HDWallet) getPkByExtendKeyAndPath(key *hdkeychain.ExtendedKey, path string) (*ecdsa.PrivateKey, error) {
	derivationPath, err := accounts.ParseDerivationPath(path)
	if err != nil {
		return nil, err
	}
	for _, n := range derivationPath {
		key, err = key.Derive(n)
		if err != nil {
			return nil, err
		}
	}
	privateKey, err := key.ECPrivKey()
	if err != nil {
		return nil, err
	}
	return (*ecdsa.PrivateKey)(privateKey), nil
}

func bid39Seed(mnemonic, password string) ([]byte, error) {
	return bip39.NewSeedWithErrorChecking(mnemonic, password)
}
