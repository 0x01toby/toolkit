package abi

import (
	"encoding/hex"
	"golang.org/x/crypto/sha3"
	"testing"
)

// keccak256
func TestSha3(t *testing.T) {
	keccak256 := sha3.NewLegacyKeccak256()
	keccak256.Write([]byte(""))
	bytes := keccak256.Sum(nil)
	toString := hex.EncodeToString(bytes)
	t.Log(toString)
}

// 标准sha3
func TestSha3_standard(t *testing.T) {
	new256 := sha3.New256()
	new256.Write([]byte(""))
	t.Log(hex.EncodeToString(new256.Sum(nil)))
}
