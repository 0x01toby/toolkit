package block

import (
	"encoding/json"
	"fmt"
	"math/big"
)

type BigInt big.Int

func (b BigInt) String() string {
	v := big.Int(b)
	return v.String()
}

func (b *BigInt) Scan(value interface{}) error {
	var t big.Int
	switch v := value.(type) {
	case []byte:
		t.SetBytes(v)
		*b = BigInt(t)
		return nil
	case string:
		t.SetString(v, 10)
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
