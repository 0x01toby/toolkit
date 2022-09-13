package block

import (
	"encoding/hex"
	"github.com/ethereum/go-ethereum/common"
	"strconv"
	"strings"
)

type Address [20]byte

func (a *Address) UnmarshalText(b []byte) error {
	return unmarshalTextByte(a[:], b, 20)
}

func (a Address) MarshalText() ([]byte, error) {
	return []byte(a.String()), nil
}

func (a Address) String() string {
	address := strings.ToLower(hex.EncodeToString(a[:]))
	hash := hex.EncodeToString(Keccak256([]byte(address)))
	ret := "0x"
	for i := 0; i < len(address); i++ {
		character := string(address[i])
		num, _ := strconv.ParseInt(string(hash[i]), 16, 64)
		if num > 7 {
			ret += strings.ToUpper(character)
		} else {
			ret += character
		}
	}
	return ret
}

func (a Address) ToCommonAddress() *common.Address {
	address := common.BytesToAddress(a[:])
	return &address
}

func Hexstr2Address(str string) Address {
	toHex := HexstrToHex(str)
	var addr Address
	copy(addr[:], toHex[len(toHex)-20:])
	return addr
}
