package block

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type Hex []byte

func (h Hex) Hex() string {
	return h.String()
}

func (h Hex) String() string {
	return hexutil.Encode(h[:])
}

func (h Hex) Bytes() []byte {
	return h[:]
}

func (h Hex) No0xPrefix() string {
	return h.String()[2:]
}

func (h Hex) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.String())
}

func (h *Hex) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		return nil
	}
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	return h.ToHex(s)
}

func (h *Hex) ToHex(value any) (err error) {
	switch v := value.(type) {
	case []byte:
		*h = v
		return
	case string:
		if len(v) > 2 {
			*h, err = hexutil.Decode(v)
		}
		return
	default:
		return fmt.Errorf("can not convert %T to HexStr", value)
	}
}

func HexstrToHex(str string) Hex {
	var h Hex
	_ = h.ToHex(str)
	return h
}

func HexstrToString(value string) string {
	var s string
	if len(value) > 64*2+2 {
		s = value[2+128:]
	} else {
		s = value[2:]
	}
	v, _ := hexutil.Decode("0x" + s)
	return string(stripBytes(v))
}

func stripBytes(input []byte) []byte {
	var result []byte
	for _, v := range input {
		if v != 0 {
			result = append(result, v)
		}
	}
	return result
}
