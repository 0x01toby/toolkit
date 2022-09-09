package rlp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/big"
	"sync"
)

var bufPool = sync.Pool{
	New: func() interface{} {
		buf := make([]byte, 8)
		return &buf
	},
}

type Cache struct {
	buf [8]byte
}

type Value struct {
	t Type
	a []*Value
	b []byte
	l uint64
	i uint64
}

var (
	valueArrayNull = &Value{t: TypeArrayNull, l: 1}
	valueNull      = &Value{t: TypeNull, l: 1}
	valueFalse     = valueNull
	valueTrue      = &Value{t: TypeBytes, b: []byte{0x1}, l: 1}
)

func (v *Value) GetString() (string, error) {
	if v.t != TypeBytes {
		return "", errNoBytes()
	}
	return string(v.b), nil
}

func (v *Value) GetElems() ([]*Value, error) {
	if v.t != TypeArray {
		return nil, errNoArray()
	}
	return v.a, nil
}

func (v *Value) GetBigInt(b *big.Int) error {
	if v.t != TypeBytes {
		return errNoBytes()
	}
	b.SetBytes(v.b)
	return nil
}

func (v *Value) GetBool() (bool, error) {
	if v.t != TypeBytes {
		return false, errNoBytes()
	}
	if bytes.Equal(v.b, valueTrue.b) {
		return true, nil
	}
	if bytes.Equal(v.b, valueFalse.b) {
		return false, nil
	}
	return false, fmt.Errorf("not a valid bool")
}

func (v *Value) Raw() []byte {
	return v.b
}

func (v *Value) Bytes() ([]byte, error) {
	if v.t != TypeBytes {
		return nil, errNoBytes()
	}
	return v.b, nil
}

func (v *Value) GetBytes(dst []byte, bits ...int) ([]byte, error) {
	if v.t != TypeBytes {
		return nil, errNoBytes()
	}
	if len(bits) > 0 {
		if len(v.b) != bits[0] {
			return nil, fmt.Errorf("bad length, expected %d but found %d", bits[0], len(v.b))
		}
	}
	dst = append(dst[:0], v.b...)
	return dst, nil
}

func (v *Value) GetAddr(buf []byte) error {
	_, err := v.GetBytes(buf, 20)
	return err
}

func (v *Value) GetHash(buf []byte) error {
	_, err := v.GetBytes(buf, 32)
	return err
}

func (v *Value) GetByte() (byte, error) {
	if v.t != TypeBytes {
		return 0, errNoBytes()
	}
	if len(v.b) != 1 {
		return 0, fmt.Errorf("bad length, expected 1 but found %d", len(v.b))
	}
	return v.b[0], nil
}

func (v *Value) GetUint64() (uint64, error) {
	if v.t != TypeBytes {
		return 0, errNoBytes()
	}
	if len(v.b) > 8 {
		return 0, fmt.Errorf("bytes %d too long for uint64", len(v.b))
	}
	buf := bufPool.Get().(*[]byte)
	defer bufPool.Put(buf)
	num := readUint(v.b, *buf)
	return num, nil
}

func (v *Value) Type() Type {
	return v.t
}

func (v *Value) Get(i int) *Value {
	if i > len(v.a) {
		return nil
	}
	return v.a[i]
}

func (v *Value) Elems() int {
	return len(v.a)
}

func (v *Value) Len() uint64 {
	if v.t == TypeArray {
		return v.l + intsize(v.l)
	}
	return v.l
}

func (v *Value) fullLen() uint64 {
	if v.t == TypeNull || v.t == TypeArrayNull {
		return 1
	}

	size := v.l
	if v.t == TypeBytes {
		if size == 1 && v.b[0] < 0x7F {
			return 1
		} else if size < 56 {
			return 1 + size
		} else {
			return 1 + intsize(size) + size
		}
	}
	if size < 56 {
		return 1 + size
	}
	return 1 + intsize(size) + size
}

func (v *Value) Set(vv *Value) {
	if v == nil || v.t != TypeArray {
		return
	}
	v.l += vv.fullLen()
	v.a = append(v.a, vv)
}

func (v *Value) marshalSize(dst []byte, short, long byte) []byte {
	if v.l < 56 {
		return append(dst, short+byte(v.l))
	}
	intSize := intsize(v.l)

	buf := bufPool.Get().(*[]byte)
	defer bufPool.Put(buf)
	binary.BigEndian.PutUint64((*buf)[:], uint64(v.l))

	dst = append(dst, long+byte(intSize))
	dst = append(dst, (*buf)[8-intSize:]...)
	return dst
}

type Type int

const (
	TypeArray Type = iota
	TypeBytes
	TypeNull
	TypeArrayNull
)

func (t Type) String() string {
	switch t {
	case TypeArray:
		return "array"
	case TypeBytes:
		return "bytes"
	case TypeNull:
		return "null"
	case TypeArrayNull:
		return "null-array"
	default:
		panic(fmt.Errorf("unknow vaue type: %d", t))

	}
}

func errNoBytes() error {
	return fmt.Errorf("value is not of type bytes")
}

func errNoArray() error {
	return fmt.Errorf("value is not of type array")
}

func readUint(b []byte, buf []byte) uint64 {
	size := len(b)
	ini := 8 - size
	for i := 0; i < ini; i++ {
		buf[i] = 0
	}
	copy(buf[ini:], b[:size])
	return binary.BigEndian.Uint64(buf[:])
}

func intsize(val uint64) uint64 {
	switch {
	case val < (1 << 8):
		return 1
	case val < (1 << 16):
		return 2
	case val < (1 << 24):
		return 3
	case val < (1 << 32):
		return 4
	case val < (1 << 40):
		return 5
	case val < (1 << 48):
		return 6
	case val < (1 << 56):
		return 7
	}
	return 8
}
