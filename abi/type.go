package abi

import (
	"fmt"
	"github.com/taorzhang/toolkit/types/block"
	"math/big"
	"reflect"
	"strings"
)

var (
	boolT         = reflect.TypeOf(false)
	uint8T        = reflect.TypeOf(uint8(0))
	uint16T       = reflect.TypeOf(uint16(0))
	uint32T       = reflect.TypeOf(uint32(0))
	uint64T       = reflect.TypeOf(uint64(0))
	int8T         = reflect.TypeOf(int8(0))
	int16T        = reflect.TypeOf(int16(0))
	int32T        = reflect.TypeOf(int32(0))
	int64T        = reflect.TypeOf(int64(0))
	addressT      = reflect.TypeOf(block.Address{})
	stringT       = reflect.TypeOf("")
	dynamicBytesT = reflect.SliceOf(reflect.TypeOf(byte(0)))
	functionT     = reflect.ArrayOf(24, reflect.TypeOf(byte(0)))
	tupleT        = reflect.TypeOf(map[string]interface{}{})
	bigIntT       = reflect.TypeOf(new(big.Int))
)

type Kind int

const (
	// KindBool is a boolean
	KindBool Kind = iota

	// KindUInt is an uint
	KindUInt

	// KindInt is an int
	KindInt

	// KindString is a string
	KindString

	// KindArray is an array
	KindArray

	// KindSlice is a slice
	KindSlice

	// KindAddress is an address
	KindAddress

	// KindBytes is a bytes array
	KindBytes

	// KindFixedBytes is a fixed bytes
	KindFixedBytes

	// KindFixedPoint is a fixed point
	KindFixedPoint

	// KindTuple is a tuple
	KindTuple

	// KindFunction is a function
	KindFunction
)

func (k Kind) String() string {
	names := [...]string{
		"Bool",
		"Uint",
		"Int",
		"String",
		"Array",
		"Slice",
		"Address",
		"Bytes",
		"FixedBytes",
		"FixedPoint",
		"Tuple",
		"Function",
	}
	return names[k]
}

type Type struct {
	kind  Kind
	size  int
	elem  *Type
	tuple []*TupleElem
	t     reflect.Type
}

func NewTupleType(inputs []*TupleElem) *Type {
	return &Type{
		kind:  KindTuple,
		tuple: inputs,
		t:     tupleT,
	}
}

func NewType(s string) (*Type, error) {
	l := newLexer(s)
	l.nextToken()
	return readType(l)

}

func (t *Type) String() string {
	return t.Format(false)
}

// TupleElems returns the elems of the tuple
func (t *Type) TupleElems() []*TupleElem {
	return t.tuple
}

func (t *Type) isVariableInput() bool {
	return t.kind == KindSlice || t.kind == KindBytes || t.kind == KindString
}

func (t *Type) isDynamicType() bool {
	if t.kind == KindTuple {
		for _, elem := range t.tuple {
			if elem.Elem.isDynamicType() {
				return true
			}
		}
		return false
	}
	return t.kind == KindString || t.kind == KindBytes || t.kind == KindSlice || (t.kind == KindArray && t.elem.isDynamicType())
}

// Decode decodes an object using this type
func (t *Type) Decode(input []byte) (interface{}, error) {
	return Decode(t, input)
}

// DecodeStruct decodes an object using this type to the out param
func (t *Type) DecodeStruct(input []byte, out interface{}) error {
	return DecodeStruct(t, input, out)
}

// Encode encodes an object using this type
func (t *Type) Encode(v interface{}) ([]byte, error) {
	return Encode(v, t)
}

// Format returns the raw representation of the type
func (t *Type) Format(includeArgs bool) string {
	switch t.kind {
	case KindTuple:
		rawAux := make([]string, 0)
		for _, i := range t.TupleElems() {
			name := i.Elem.Format(includeArgs)
			if i.Indexed {
				name += " indexed"
			}
			if includeArgs {
				if i.Name != "" {
					name += " " + i.Name
				}
			}
			rawAux = append(rawAux, name)
		}
		return fmt.Sprintf("tuple(%s)", strings.Join(rawAux, ","))

	case KindArray:
		return fmt.Sprintf("%s[%d]", t.elem.Format(includeArgs), t.size)

	case KindSlice:
		return fmt.Sprintf("%s[]", t.elem.Format(includeArgs))

	case KindBytes:
		return "bytes"

	case KindFixedBytes:
		return fmt.Sprintf("bytes%d", t.size)

	case KindString:
		return "string"

	case KindBool:
		return "bool"

	case KindAddress:
		return "address"

	case KindFunction:
		return "function"

	case KindUInt:
		return fmt.Sprintf("uint%d", t.size)

	case KindInt:
		return fmt.Sprintf("int%d", t.size)

	default:
		panic(fmt.Errorf("BUG: abi type not found %s", t.kind.String()))
	}
}

// TupleElem is an element of a tuple
type TupleElem struct {
	Name    string
	Elem    *Type
	Indexed bool
}

func getTypeSize(t *Type) int {
	if t.kind == KindArray && !t.elem.isDynamicType() {
		if t.elem.kind == KindArray || t.elem.kind == KindTuple {
			return t.size * getTypeSize(t.elem)
		}
		return t.size * 32
	} else if t.kind == KindTuple && !t.elem.isDynamicType() {
		total := 0
		for _, elem := range t.tuple {
			total += getTypeSize(elem.Elem)
		}
		return total
	}
	return 32
}

type ArgumentStr struct {
	Name       string
	Type       string
	Indexed    bool
	Components []*ArgumentStr
}

func NewTupleTypeFromArgs(inputs []*ArgumentStr) (*Type, error) {
	elems := make([]*TupleElem, 0)
	for _, i := range inputs {
		typ, err := NewTypeFromArgument(i)
		if err != nil {
			return nil, err
		}
		elems = append(elems, &TupleElem{
			Name:    i.Name,
			Elem:    typ,
			Indexed: i.Indexed,
		})
	}
	return NewTupleType(elems), nil
}

func NewTypeFromArgument(arg *ArgumentStr) (*Type, error) {
	str, err := parseType(arg)
	if err != nil {
		return nil, err
	}
	return NewType(str)
}

func parseType(arg *ArgumentStr) (string, error) {
	if !strings.HasPrefix(arg.Type, "tuple") {
		return arg.Type, nil
	}

	if len(arg.Components) == 0 {
		return "tuple()", nil
	}

	// parse the arg components from the tuple
	str := make([]string, 0)
	for _, i := range arg.Components {
		aux, err := parseType(i)
		if err != nil {
			return "", err
		}
		if i.Indexed {
			str = append(str, aux+" indexed "+i.Name)
		} else {
			str = append(str, aux+" "+i.Name)
		}
	}
	return fmt.Sprintf("tuple(%s)%s", strings.Join(str, ","), strings.TrimPrefix(arg.Type, "tuple")), nil
}
