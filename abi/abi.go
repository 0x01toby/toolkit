package abi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/taorzhang/toolkit/types/block"
	"io"
	"strings"
)

// ABI 一个合约的ABI
type ABI struct {
	Constructor *Method
	// MethodsBySignature method signature => Method
	MethodsBySignature map[string]*Method
	// Events event signature => Event
	EventsBySignature map[string]*Event
	Errors            map[string]*Error
}

func NewABI(s string) (*ABI, error) {
	return NewABIFromReader(bytes.NewReader([]byte(s)))
}

func NewABIFromList(humanReadableAbi []string) (*ABI, error) {
	abi := newPlainABI()
	for _, c := range humanReadableAbi {
		if strings.HasPrefix(c, "constructor") {
			typ, err := NewType("tuple" + strings.TrimPrefix(c, "constructor"))
			if err != nil {
				return nil, err
			}
			abi.Constructor = &Method{
				Inputs: typ,
			}
		} else if strings.HasPrefix(c, "function ") {
			method, err := NewMethod(c)
			if err != nil {
				return nil, err
			}
			abi.addMethod(method)
		} else if strings.HasPrefix(c, "event ") {
			event, err := NewEvent(c)
			if err != nil {
				return nil, err
			}
			abi.addEvent(event)
		} else if strings.HasPrefix(c, "error ") {
			errTyp, err := NewError(c)
			if err != nil {
				return nil, err
			}
			abi.addError(errTyp)
		} else {
			return nil, fmt.Errorf("either event or function expected")
		}
	}
	return abi, nil
}

func NewABIFromReader(r io.Reader) (*ABI, error) {
	abi := newPlainABI()
	abi.Constructor = nil
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&abi); err != nil {
		return nil, err
	}
	return abi, nil
}

func newPlainABI() *ABI {
	return &ABI{
		Constructor:        new(Method),
		MethodsBySignature: make(map[string]*Method),
		EventsBySignature:  make(map[string]*Event),
		Errors:             make(map[string]*Error),
	}
}

func (a *ABI) GetMethodBySig(sig string) *Method {
	return a.MethodsBySignature[sig]
}

func (a *ABI) GetMethodByID(id string) *Method {
	for _, method := range a.MethodsBySignature {
		if strings.EqualFold(method.HexID(), id) {
			return method
		}
	}
	return nil
}

func (a *ABI) GetEventBySig(sig string) *Event {
	return a.EventsBySignature[sig]
}

// GetEventByID 根据ID获取Event
func (a *ABI) GetEventByID(id string) *Event {
	for _, event := range a.EventsBySignature {
		if strings.EqualFold(event.ID().String(), id) {
			return event
		}
	}
	return nil
}

func (a *ABI) addError(e *Error) {
	a.Errors[e.Name] = e
}

func (a *ABI) addEvent(e *Event) {
	a.EventsBySignature[e.Sig()] = e
}

func (a *ABI) addMethod(m *Method) {
	a.MethodsBySignature[m.Sig()] = m
}

func (a *ABI) UnmarshalJSON(data []byte) error {
	var fields []struct {
		Type            string
		Name            string
		Constant        bool
		Anonymous       bool
		StateMutability string
		Inputs          []*ArgumentStr
		Outputs         []*ArgumentStr
	}

	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}

	for _, field := range fields {
		switch field.Type {
		case "constructor":
			if a.Constructor != nil {
				return fmt.Errorf("multiple constructor declaration")
			}
			input, err := NewTupleTypeFromArgs(field.Inputs)
			if err != nil {
				panic(err)
			}
			a.Constructor = &Method{
				Inputs: input,
			}

		case "function", "":
			c := field.Constant
			if field.StateMutability == "view" || field.StateMutability == "pure" {
				c = true
			}

			inputs, err := NewTupleTypeFromArgs(field.Inputs)
			if err != nil {
				panic(err)
			}
			outputs, err := NewTupleTypeFromArgs(field.Outputs)
			if err != nil {
				panic(err)
			}
			method := &Method{
				Name:    field.Name,
				Const:   c,
				Inputs:  inputs,
				Outputs: outputs,
			}
			a.addMethod(method)

		case "event":
			input, err := NewTupleTypeFromArgs(field.Inputs)
			if err != nil {
				panic(err)
			}
			event := &Event{
				Name:      field.Name,
				Anonymous: field.Anonymous,
				Inputs:    input,
			}
			a.addEvent(event)

		case "error":
			input, err := NewTupleTypeFromArgs(field.Inputs)
			if err != nil {
				panic(err)
			}
			errObj := &Error{
				Name:   field.Name,
				Inputs: input,
			}
			a.addError(errObj)

		case "fallback":
		case "receive":
			// do nothing

		default:
			return fmt.Errorf("unknown field type '%s'", field.Type)
		}
	}
	return nil
}

// Method abi method
type Method struct {
	// Name 方法名称
	// example function transfer(address,address,uint256) public returns(bool) name就是transfer
	Name  string
	Const bool
	// Inputs 入参
	// example function transfer(address,address,uint256) public returns(bool), inputs就是tuple(address,address,uint256)
	Inputs *Type
	// Outputs 出参
	// example function transfer(address,address,uint256)public returns(bool) outputs就是tuple(bool)
	Outputs *Type
}

func NewMethod(name string) (*Method, error) {
	name, inputs, outputs, err := parseMethodSignature(name)
	if err != nil {
		return nil, err
	}
	m := &Method{Name: name, Inputs: inputs, Outputs: outputs}
	return m, nil
}

// Sig 构造函数签名（没有返回值）
// example function transfer(address from, address to, uint256 amount) public returns(bool)
// 返回是：transfer(address,address,uint256)
func (m *Method) Sig() string {
	return buildSignature(m.Name, m.Inputs)
}

// ID 根据Sig 然后通过keccak256 hash后的结果签名4bytes
// example keccak256("transfer(address,address,uint256")[0:4]
// 4个byte 32位 4位是一个16进制字符 -> 8个16进制字符
func (m *Method) ID() []byte {
	k := acqKeccak()
	defer releaseKeccak(k)
	k.Write([]byte(m.Sig()))
	dst := k.Sum(nil)[:4]
	return dst
}

func (m *Method) HexID() string {
	return hexutil.Encode(m.ID())
}

// Encode encode一个method + args 到call data
func (m *Method) Encode(args interface{}) ([]byte, error) {
	data, err := Encode(args, m.Inputs)
	if err != nil {
		return nil, err
	}
	data = append(m.ID(), data...)
	return data, nil
}

// Decode decode call data 到 method return outputs
func (m *Method) Decode(data []byte) (map[string]interface{}, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty response")
	}
	resp, err := Decode(m.Outputs, data)
	if err != nil {
		return nil, err
	}
	return resp.(map[string]interface{}), nil
}

type Event struct {
	Name      string
	Anonymous bool
	Inputs    *Type
}

func NewEvent(name string) (*Event, error) {
	name, typ, err := parseEventSignature(name)
	if err != nil {
		return nil, err
	}
	return &Event{Name: name, Inputs: typ}, nil
}

// Sig 构造函数名称
func (e *Event) Sig() string {
	return buildSignature(e.Name, e.Inputs)
}

// ID 根据函数的signature通过keccek256 hash得到256位bytes
// example keccak256("event transfer(address,address,uint256)") 得到256bytes
func (e *Event) ID() (id block.Hash) {
	k := acqKeccak()
	defer releaseKeccak(k)
	k.Write([]byte(e.Sig()))
	dst := k.Sum(nil)
	copy(id[:], dst)
	return
}

// ParseLog 解析日志
func (e *Event) ParseLog(log *block.Log) (map[string]interface{}, error) {
	if !e.Match(log) {
		return nil, fmt.Errorf("log does not match this event")
	}
	return ParseLog(e.Inputs, log)
}

// Match 判断topic0是否匹配
func (e *Event) Match(log *block.Log) bool {
	if len(log.Topics) == 0 {
		return false
	}
	return log.Topics[0] == e.ID()
}

type Error struct {
	Name   string
	Inputs *Type
}

func NewError(name string) (*Error, error) {
	name, typ, err := parseErrorSignature(name)
	if err != nil {
		return nil, err
	}
	return &Error{Name: name, Inputs: typ}, nil
}
