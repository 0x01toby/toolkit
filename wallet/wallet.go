package wallet

import (
	"context"
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/taorzhang/toolkit/abi"
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

// GetNonce 获取交易nonce
func (w *Wallet) GetNonce(status NonceStatus) (nonce uint64, err error) {
	nonce, err = w.Client.GetNonce(context.Background(), block.Hexstr2Address(w.Address()), status.ToString())
	return
}

// EstimateGas 评估gas 费用
func (w *Wallet) EstimateGas(txData *types.Transaction) (*big.Int, error) {
	parameter := client.CallParameter{
		From:     *block.Hexstr2Address(w.Address()).ToCommonAddress(),
		Data:     block.Hex(txData.Data()),
		GasPrice: txData.GasPrice(),
		Value:    txData.Value(),
	}
	if txData.To() != nil {
		parameter.To = block.Hexstr2Address(txData.To().Hex()).ToCommonAddress()
	}
	return w.Client.EstimateGas(context.Background(), parameter)
}

// DeployContract 部署合约
// code none 0x prefix
func (w *Wallet) DeployContract(code string) (*block.Hash, error) {
	txData, err := w.createContractTxData(code, "pending")
	if err != nil {
		return nil, err
	}
	sigedTx, err := w.signTx(txData)
	if err != nil {
		return nil, err
	}
	sendTx, err := w.Client.SendTx(context.Background(), sigedTx)
	if err != nil {
		return nil, err
	}
	address := block.Hex2Hash(sendTx)
	return &address, nil
}

// SendNativeToken 发送原生代币
func (w *Wallet) SendNativeToken(to block.Address, amount *big.Int) (*block.Hash, error) {
	txData, err := w.createLegacyTxData(Pending, to, amount, nil)
	if err != nil {
		return nil, err
	}
	tx, err := w.signTx(txData)
	if err != nil {
		return nil, err
	}
	txHash, err := w.Client.SendTx(context.Background(), tx)
	if err != nil {
		return nil, err
	}
	hash := block.Hex2Hash(txHash)
	return &hash, nil
}

// SendErc20Token 发送erc20代币
func (w *Wallet) SendErc20Token(contract block.Address, to block.Address, amount string) (*block.Hash, error) {
	method, err := abi.NewMethod("function transfer(address dst, uint256 wad)")
	if err != nil {
		return nil, err
	}
	decimals, err := w.getContractDecimals(contract)
	if err != nil {
		return nil, err
	}
	tokenNum, err := DataMulDecimal(amount, int(decimals))
	if err != nil {
		return nil, err
	}
	txData, err := w.createLegacyTxData(Pending, contract, big.NewInt(0), method, to.String(), tokenNum.String())
	if err != nil {
		return nil, err
	}
	tx, err := w.signTx(txData)
	if err != nil {
		return nil, err
	}
	sendTx, err := w.Client.SendTx(context.Background(), tx)
	if err != nil {
		return nil, err
	}
	hash := block.Hex2Hash(sendTx)
	return &hash, err
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

// signTx 对交易进行签名
func (w *Wallet) signTx(tx types.TxData) (string, error) {
	chainID, err := w.Client.ChainID(context.Background())
	if err != nil {
		return "", err
	}
	signTx, err := types.SignTx(types.NewTx(tx), types.NewEIP155Signer(chainID), w.PrivateKey)
	if err != nil {
		return "", err
	}
	binary, err := signTx.MarshalBinary()
	if err != nil {
		return "", err
	}
	return hexutil.Encode(binary), nil
}

// createContractTxData 创建合约交易
func (w *Wallet) createContractTxData(code string, nonceStatus NonceStatus) (types.TxData, error) {
	nonce, err := w.GetNonce(nonceStatus)
	if err != nil {
		return nil, err
	}
	gasPrice, err := w.Client.GetGasPrice(context.Background())
	if err != nil {
		return nil, err
	}
	txData := &types.LegacyTx{
		Nonce:    nonce,
		GasPrice: gasPrice,
		Value:    big.NewInt(0),
		Data:     common.FromHex(code),
	}
	gas, err := w.EstimateGas(types.NewTx(txData))
	if err != nil {
		return nil, err
	}
	txData.Gas = gas.Uint64()
	return txData, nil
}

// createLegacyTxData 创建一笔legacy交易
func (w *Wallet) createLegacyTxData(nonceStatus NonceStatus, to block.Address, amount *big.Int, method *abi.Method, args ...interface{}) (types.TxData, error) {
	nonce, err := w.GetNonce(nonceStatus)
	if err != nil {
		return nil, err
	}
	gasPrice, err := w.Client.GetGasPrice(context.Background())
	if err != nil {
		return nil, err
	}
	var encode string
	if method != nil {
		result, err := method.Encode(args)
		if err != nil {
			return nil, err
		}
		encode = block.Hex(result).Hex()
	}
	txData := &types.LegacyTx{
		Nonce:    nonce,
		GasPrice: gasPrice,
		Value:    amount,
		To:       to.ToCommonAddress(),
		Data:     common.FromHex(encode),
	}
	gas, err := w.EstimateGas(types.NewTx(txData))
	if err != nil {
		return nil, err
	}
	txData.Gas = gas.Uint64()
	return txData, nil
}

func (w *Wallet) getContractDecimals(contractAddress block.Address) (uint64, error) {
	var decimals string
	method, err := abi.NewMethod("function decimals()")
	if err != nil {
		return 0, err
	}
	to := block.Hexstr2Address("0x11fE4B6AE13d2a6055C8D9cF65c55bac32B5d844")
	err = w.Client.MethodCall(context.Background(), &decimals, client.CallParameter{
		To:   to.ToCommonAddress(),
		Data: common.FromHex(block.Hex(method.ID()).Hex()),
		Gas:  client.DefaultGasLimitInt,
	}.ToArg(), "latest")
	bigInt, err := block.HexStrToBigInt(decimals)
	if err != nil {
		return 0, err
	}
	return bigInt.Uint64(), nil
}
