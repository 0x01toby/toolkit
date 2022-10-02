package wallet

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

type PrivateKey struct {
	*ecdsa.PrivateKey
}

func NewPkFromPkStr(pkStr string) (*PrivateKey, error) {
	key, err := privateKeyStr2privateKey(pkStr)
	if err != nil {
		return nil, err
	}
	return &PrivateKey{PrivateKey: key}, nil
}

func NewPkFromKeyStore(keyStoreContent []byte, password string) (*PrivateKey, error) {
	key, err := keyStoreToPrivateKey(keyStoreContent, password)
	if err != nil {
		return nil, err
	}
	return &PrivateKey{PrivateKey: key}, nil
}

func (pk *PrivateKey) PublicAddress() string {
	return crypto.PubkeyToAddress(pk.PublicKey).Hex()
}

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

func keyStoreToPrivateKey(keyStoreContent []byte, password string) (*ecdsa.PrivateKey, error) {
	unlockedKey, err := keystore.DecryptKey(keyStoreContent, password)
	if err != nil {
		return nil, err
	}
	return unlockedKey.PrivateKey, nil
}
