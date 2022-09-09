package block

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
)

type BigInt big.Int

func (b BigInt) String() string {
	v := big.Int(b)
	return v.String()
}

func (b BigInt) ToBigInt() *big.Int {
	b2 := big.Int(b)
	return &b2
}

func (b *BigInt) Scan(value interface{}) error {
	var t big.Int
	switch v := value.(type) {
	case []byte:
		t.SetBytes(v)
		*b = BigInt(t)
		return nil
	case string:
		if strings.HasPrefix(strings.ToLower(v), "0x") {
			t.SetString(v[2:], 16)
		} else {
			t.SetString(v, 10)
		}
		*b = BigInt(t)
		return nil
	}
	return fmt.Errorf("can't convert %T to hex.BigInt", value)
}

func (b BigInt) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.String())
}

func (b *BigInt) UnmarshalJSON(value []byte) error {
	var s string
	if err := json.Unmarshal(value, &s); err != nil {
		return err
	}
	return b.Scan(s)
}

func HexStrToBigInt(str string) (*big.Int, error) {
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
