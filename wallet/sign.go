package wallet

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
)

func SignTransaction(chainID int64, tx *types.Transaction, privateKey *ecdsa.PrivateKey) (string, error) {
	signTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(chainID)), privateKey)
	if err != nil {
		return "", err
	}
	binary, err := signTx.MarshalBinary()
	if err != nil {
		return "", err
	}
	return hexutil.Encode(binary), nil
}
