package block

import "encoding/hex"

type Hash [32]byte

func Hex2Hash(str string) Hash {
	h := Hash{}
	_ = h.UnmarshalText(completeHex(str, 32))
	return h
}

func (h *Hash) UnmarshalText(b []byte) error {
	return unmarshalTextByte(h[:], b, 32)
}

func (h Hash) MarshalText() ([]byte, error) {
	return []byte(h.String()), nil
}

func (h Hash) String() string {
	return "0x" + hex.EncodeToString(h[:])
}

func (h Hash) Location() string {
	return h.String()
}

func (h Hash) Bytes() []byte {
	return h[:]
}
