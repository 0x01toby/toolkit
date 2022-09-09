package wallet

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/taorzhang/toolkit/abi"
	"github.com/taorzhang/toolkit/types/block"
	"math/big"
	"strings"
)

func CreateLegacyTx(
	nonce uint64,
	to block.Address,
	amount *big.Int,
	gasLimit uint64,
	gasPrice *big.Int,
	method *abi.Method,
	args ...interface{}) (types.TxData, error) {
	var encode []byte
	if method != nil {
		result, err := method.Encode(args)
		if err != nil {
			return nil, err
		}
		copy(encode[:], result[:])
	}
	return &types.LegacyTx{
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      gasLimit,
		Value:    amount,
		To:       to.ToCommonAddress(),
		Data:     encode,
	}, nil
}

func Create1559Tx(
	chainID *big.Int,
	nonce uint64,
	to block.Address,
	amount *big.Int,
	gasLimit uint64,
	gasTipCap *big.Int,
	gasFeeCap *big.Int,
	accessList types.AccessList,
	method *abi.Method,
	args ...interface{}) (types.TxData, error) {
	var encode []byte
	if method != nil {
		result, err := method.Encode(args)
		if err != nil {
			return nil, err
		}
		copy(encode[:], result[:])
	}
	return &types.DynamicFeeTx{
		ChainID:    chainID,
		Nonce:      nonce,
		To:         to.ToCommonAddress(),
		V:          amount,
		Gas:        gasLimit,
		GasTipCap:  gasTipCap,
		GasFeeCap:  gasFeeCap,
		AccessList: accessList,
		Data:       encode,
	}, nil
}

func CreateContract(nonce uint64, gasLimit uint64, gasPrice *big.Int, code string) (types.TxData, error) {
	return &types.LegacyTx{
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      gasLimit,
		Value:    big.NewInt(0),
		Data:     common.FromHex(code),
	}, nil
}

func ToWei(num string) (*big.Int, error) {
	return DataMulDecimal(num, 18)
}

// dealDecimal
// input 18 output 10^18
func dealDecimal(decimals int) *big.Int {
	i := big.NewInt(int64(decimals))
	exp := i.Exp(big.NewInt(10), i, nil)
	return exp
}

func DataMulDecimal(num string, decimal int) (*big.Int, error) {
	var data big.Float
	setString, b := data.SetString(num)
	if !b {
		return nil, fmt.Errorf("convert '%s' to big.float failed", num)
	}
	var decimalFloat big.Float
	dFloat := decimalFloat.SetInt(dealDecimal(decimal))
	mul := setString.Mul(setString, dFloat)
	text := mul.Text('f', 10)
	split := strings.Split(text, ".")
	var result big.Int
	i, b2 := result.SetString(split[0], 10)
	if !b2 {
		return nil, fmt.Errorf("convert '%s' to big.Int failed", split[0])
	}
	return i, nil
}
