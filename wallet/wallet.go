package wallet

import (
	"context"
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/taorzhang/toolkit/client"
	"github.com/taorzhang/toolkit/types/block"
	"math/big"
)

type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	Client     client.Provider
}

func NewWalletFromPrivateKeyStr(privateKeyString string, options ...Option) *Wallet {
	key, err := privateKeyStr2privateKey(privateKeyString)
	if err != nil {
		panic(err)
	}
	wallet := &Wallet{PrivateKey: key}
	for _, opt := range options {
		opt(wallet)
	}
	return wallet
}

func NewWalletFromKeyStore(keyStoreContent []byte, password string, options ...Option) *Wallet {
	key, err := KeyStoreToPrivateKey(keyStoreContent, password)
	if err != nil {
		panic(err)
	}
	wallet := &Wallet{PrivateKey: key}
	for _, opt := range options {
		opt(wallet)
	}
	return wallet
}

// Address 根据私钥生成钱包地址
func (w *Wallet) Address() string {
	return crypto.PubkeyToAddress(w.PrivateKey.PublicKey).Hex()
}

func (w *Wallet) GetNonce(status string) (nonce uint64, err error) {
	nonce, err = w.Client.GetNonce(context.Background(), block.Hexstr2Address(w.Address()), status)
	return
}

func (w *Wallet) EstimateGas(txData *types.Transaction) {
	w.Client.EstimateGas(context.Background(), client.CallParameter{
		From:     w.Address(),
		To:       txData.To().Hex(),
		Data:     block.Hex(txData.Data()).Hex(),
		Gas:      big.NewInt(int64(txData.Gas())).String(),
		GasPrice: txData.GasPrice().String(),
		Value:    txData.Value().String(),
	})
}

// CreateContract 部署合约
func (w *Wallet) CreateContract() {

}

// SendNativeToken 发送原生代币
func (w *Wallet) SendNativeToken(amount *big.Int) {

}

// SendErc20Token 发送erc20代币
func (w *Wallet) SendErc20Token(contract block.Address, amount *big.Int) {

}

// ApprovalErc20Token 授权erc20代币
func (w *Wallet) ApprovalErc20Token(contract block.Address, amount *big.Int) {

}

// ApprovalAllErc20Token 授权所有的erc20代币
func (w *Wallet) ApprovalAllErc20Token(contract block.Address) {

}

// SendErc721 发送erc721
func (w *Wallet) SendErc721(contract block.Address, tokenID *big.Int) {

}

// ApprovalErc721 授权
func (w *Wallet) ApprovalErc721(contract block.Address, tokenID *big.Int) {

}

func (w *Wallet) SendErc1155(contract block.Address, tokenID *big.Int, amount *big.Int) {

}

func (w *Wallet) BatchSendErc1155(contract block.Address) {

}

func (w *Wallet) ApprovalErc1155(contract block.Address, tokenID *big.Int) {

}

func (w *Wallet) sigTx(txData types.TxData) ([]byte, error) {
	chainID, err := w.Client.ChainID(context.Background())
	if err != nil {
		return nil, err
	}
	signTx, err := types.SignTx(types.NewTx(txData), types.NewEIP155Signer(chainID), w.PrivateKey)
	if err != nil {
		return nil, err
	}
	bytes, err := rlp.EncodeToBytes(signTx)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
