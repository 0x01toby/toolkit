package abi

import (
	"encoding/binary"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/taorzhang/toolkit/types/block"
	"math/big"
	"reflect"
	"strconv"
)

func Decode(t *Type, input []byte) (interface{}, error) {
	if len(input) == 0 {
		return nil, fmt.Errorf("empty input")
	}
	val, _, err := decode(t, input)
	return val, err
}

func DecodeStruct(t *Type, input []byte, out interface{}) error {
	val, err := Decode(t, input)
	if err != nil {
		return err
	}

	dc := &mapstructure.DecoderConfig{
		Result:           out,
		WeaklyTypedInput: true,
		TagName:          "abi",
	}

	ms, err := mapstructure.NewDecoder(dc)
	if err != nil {
		return err
	}

	if err = ms.Decode(val); err != nil {
		return err
	}
	return nil
}

func decode(t *Type, input []byte) (interface{}, []byte, error) {
	var data []byte
	var length int
	var err error

	if len(input) < 32 {
		return nil, nil, fmt.Errorf("incorret length")
	}

	if t.isVariableInput() {
		length, err = readLength(input)
		if err != nil {
			return nil, nil, err
		}
	} else {
		data = input[:32]
	}

	switch t.kind {
	case KindTuple:
		return decodeTuple(t, input)

	case KindSlice:
		return decodeArraySlice(t, input[32:], length)

	case KindAddress:
		return decodeArraySlice(t, input, t.size)
	}

	var val interface{}

	switch t.kind {
	case KindBool:
		val, err = decodeBool(data)

	case KindInt, KindUInt:
		val = decodeInteger(t, data)

	case KindString:
		val = string(input[32 : 32+length])

	case KindBytes:
		val = input[32 : 32+length]

	case KindAddress:
		val, err = decodeAddress(data)

	case KindFixedBytes:
		val, err = decodeFixedBytes(t, data)

	case KindFunction:
		val, err = decodeFunction(t, data)

	default:
		return nil, nil, fmt.Errorf("decoding not available for type '%s'", t.kind)
	}
	return val, input[32:], err
}

func decodeFunction(t *Type, word []byte) ([24]byte, error) {
	res := [24]byte{}
	if !allZeros(word[24:32]) {
		return res, fmt.Errorf("function type expects the last 8 bytes to be empty but found: %b", word[24:32])
	}
	copy(res[:], word[0:24])
	return res, nil
}

func decodeFixedBytes(t *Type, word []byte) (interface{}, error) {
	array := reflect.New(t.t).Elem()
	reflect.Copy(array, reflect.ValueOf(word[0:t.size]))
	return array.Interface(), nil
}

func decodeAddress(b []byte) (block.Address, error) {
	var s block.Address
	if len(b) != 32 {
		return s, fmt.Errorf("len is not correct")
	}
	copy(s[:], b[12:])
	return s, nil
}

func decodeInteger(t *Type, b []byte) interface{} {
	switch t.t.Kind() {
	case reflect.Uint8:
		return b[len(b)-1]

	case reflect.Uint16:
		return binary.BigEndian.Uint16(b[len(b)-2:])

	case reflect.Uint32:
		return binary.BigEndian.Uint32(b[len(b)-4:])

	case reflect.Uint64:
		return binary.BigEndian.Uint64(b[len(b)-8:])

	case reflect.Int8:
		return int8(b[len(b)-1])

	case reflect.Int16:
		return int16(binary.BigEndian.Uint16(b[len(b)-2:]))

	case reflect.Int32:
		return int32(binary.BigEndian.Uint32(b[len(b)-4:]))

	case reflect.Int64:
		return int64(binary.BigEndian.Uint64(b[len(b)-8:]))

	default:
		ret := new(big.Int).SetBytes(b)
		if t.kind == KindUInt {
			return ret
		}
		if ret.Cmp(maxInt256) > 0 {
			ret.Add(maxUint256, big.NewInt(0).Neg(ret))
			ret.Add(ret, big.NewInt(1))
			ret.Neg(ret)
		}
		return ret
	}
}

func decodeBool(data []byte) (interface{}, error) {
	switch data[31] {
	case 0:
		return false, nil
	case 1:
		return true, nil
	default:
		return false, fmt.Errorf("bad boolean")
	}
}

func decodeArraySlice(t *Type, data []byte, size int) (interface{}, []byte, error) {
	if size < 0 {
		return nil, nil, fmt.Errorf("size is lower than zero")
	}
	if 32*size > len(data) {
		return nil, nil, fmt.Errorf("size is too big")
	}

	var res reflect.Value
	if t.kind == KindSlice {
		res = reflect.MakeSlice(t.t, size, size)
	} else if t.kind == KindArray {
		res = reflect.New(t.t).Elem()
	}

	origin := data
	originLen := len(origin)

	for idx := 0; idx < size; idx++ {
		isDynamic := t.elem.isDynamicType()

		if len(data) < 32 {
			return nil, nil, fmt.Errorf("incorrect length")
		}

		entry := data
		if isDynamic {
			offset, err := readOffset(data, originLen)
			if err != nil {
				return nil, nil, err
			}
			entry = origin[offset:]
		}

		val, tail, err := decode(t.elem, entry)
		if err != nil {
			return nil, nil, err
		}
		if !isDynamic {
			data = tail
		} else {
			data = data[32:]
		}
		res.Index(idx).Set(reflect.ValueOf(val))
	}
	return res.Interface(), data, nil
}

func decodeTuple(t *Type, data []byte) (interface{}, []byte, error) {
	res := make(map[string]interface{})

	origin := data
	originLen := len(origin)

	for idx, arg := range t.tuple {
		if len(data) < 32 {
			return nil, nil, fmt.Errorf("incorrect length")
		}
		entry := data
		if arg.Elem.isDynamicType() {
			offset, err := readOffset(data, originLen)
			if err != nil {
				return nil, nil, err
			}
			entry = origin[offset:]
		}

		val, tail, err := decode(arg.Elem, entry)
		if err != nil {
			return nil, nil, err
		}

		if !arg.Elem.isDynamicType() {
			data = tail
		} else {
			data = data[32:]
		}

		name := arg.Name
		if name == "" {
			name = strconv.Itoa(idx)
		}
		if _, ok := res[name]; !ok {
			res[name] = val
		} else {
			return nil, nil, fmt.Errorf("tuple with repeated values")
		}
	}
	return res, data, nil
}

func readLength(data []byte) (int, error) {
	lengthBig := big.NewInt(0).SetBytes(data[0:32])
	if lengthBig.BitLen() > 63 {
		return 0, fmt.Errorf("length larger than int64: %v", lengthBig.Int64())
	}
	length := int(lengthBig.Uint64())
	if length > len(data) {
		return 0, fmt.Errorf("length insufficient %v require %v", len(data), length)
	}
	return length, nil
}

func readOffset(data []byte, len int) (int, error) {
	offsetBig := big.NewInt(0).SetBytes(data[0:32])
	if offsetBig.BitLen() > 63 {
		return 0, fmt.Errorf("offset larger than int64: %v", offsetBig.Int64())
	}
	offset := int(offsetBig.Int64())
	if offset > len {
		return 0, fmt.Errorf("offset insufficient %v require %v", len, offset)
	}
	return offset, nil
}

func allZeros(b []byte) bool {
	for _, i := range b {
		if i != 0 {
			return false
		}
	}
	return true
}

var (
	maxUint256 = big.NewInt(0).Add(
		big.NewInt(0).Exp(big.NewInt(2), big.NewInt(256), nil),
		big.NewInt(-1))
	maxInt256 = big.NewInt(0).Add(
		big.NewInt(0).Exp(big.NewInt(2), big.NewInt(255), nil),
		big.NewInt(-1))
)
