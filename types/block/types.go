package block

import (
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/sha3"
	"strings"
)

type AccessEntry struct {
	Address Address `json:"address"`
	Storage []Hash  `json:"storageKeys"`
}

type AccessList []AccessEntry

func unmarshalTextByte(dst, src []byte, size int) error {
	str := string(src)

	str = strings.Trim(str, "\"")
	if !strings.HasPrefix(str, "0x") {
		return fmt.Errorf("0x prefix not found")
	}
	str = str[2:]
	b, err := hex.DecodeString(str)
	if err != nil {
		return err
	}
	if len(b) != size {
		return fmt.Errorf("length %d is not correct, expected %d", len(b), size)
	}
	copy(dst, b)
	return nil
}

func completeHex(str string, num int) []byte {
	num = num * 2
	str = strings.TrimPrefix(str, "0x")

	size := len(str)
	if size < num {
		for i := size; i < num; i++ {
			str = "0" + str
		}
	} else {
		diff := size - num
		str = str[diff:]
	}
	return []byte("0x" + str)
}

func Keccak256(v ...[]byte) []byte {
	h := sha3.NewLegacyKeccak256()
	for _, i := range v {
		h.Write(i)
	}
	return h.Sum(nil)
}
