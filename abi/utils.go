package abi

import (
	"fmt"
	"golang.org/x/crypto/sha3"
	_ "golang.org/x/crypto/sha3"
	"hash"
	"strings"
	"sync"
)

var keccakPool = sync.Pool{
	New: func() interface{} {
		return sha3.NewLegacyKeccak256()
	},
}

func acqKeccak() hash.Hash {
	return keccakPool.Get().(hash.Hash)
}

func releaseKeccak(k hash.Hash) {
	k.Reset()
	keccakPool.Put(k)
}

// buildSignature 生成函数的签名
// like transfer(address,address,uint256)
func buildSignature(name string, typ *Type) string {
	types := make([]string, len(typ.tuple))
	for i, input := range typ.tuple {
		types[i] = strings.Replace(input.Elem.String(), "tuple", "", -1)
	}
	return fmt.Sprintf("%v(%v)", name, strings.Join(types, ","))
}
