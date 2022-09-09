package wallet

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func privateKeyStr2privateKey(privateKeyStr string) (*ecdsa.PrivateKey, error) {
	privateKeyByte, err := hexutil.Decode(privateKeyStr)
	if err != nil {
		return nil, err
	}
	privateKey, err := crypto.ToECDSA(privateKeyByte)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

func KeyStoreToPrivateKey(keyStoreContent []byte, password string) (*ecdsa.PrivateKey, error) {
	unlockedKey, err := keystore.DecryptKey(keyStoreContent, password)
	if err != nil {
		return nil, err
	}
	return unlockedKey.PrivateKey, nil
}