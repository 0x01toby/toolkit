package client

import (
	"fmt"
	"math/big"
	"strings"
)

func hexStrToBigInt(str string) (*big.Int, error) {
	if !strings.HasPrefix(strings.ToLower(str), "0x") {
		return nil, fmt.Errorf("str '%s' not a hex string", str)
	}
	var b *big.Int
	setString, ok := b.SetString(str[2:], 16)
	if !ok {
		return nil, fmt.Errorf("str '%s' can not convert to big.Int", str)
	}
	return setString, nil
}
