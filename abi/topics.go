package abi

import (
	"bytes"
	"fmt"
	"github.com/taorzhang/toolkit/types/block"
)

func init() {
	topicTrue[31] = 1
}

func ParseLog(args *Type, log *block.Log) (map[string]interface{}, error) {
	var indexed, nonIndexed []*TupleElem
	for _, arg := range args.TupleElems() {
		if arg.Indexed {
			indexed = append(indexed, arg)
		} else {
			nonIndexed = append(nonIndexed, arg)
		}
	}
	indexedObjects, err := ParseTopics(&Type{kind: KindTuple, tuple: indexed}, log.Topics[1:])
	if err != nil {
		return nil, err
	}
	var nonIndexedObjects map[string]interface{}
	if len(nonIndexed) > 0 {
		nonIndexedRaw, err := Decode(&Type{kind: KindTuple, tuple: nonIndexed}, log.Data)
		if err != nil {
			return nil, err
		}
		raw, ok := nonIndexedRaw.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("bad decoding")
		}
		nonIndexedObjects = raw
	}
	res := make(map[string]interface{})
	for _, arg := range args.TupleElems() {
		{
			if arg.Indexed {
				res[arg.Name] = indexedObjects[0]
				indexedObjects = indexedObjects[1:]
			} else {
				res[arg.Name] = nonIndexedObjects[arg.Name]
			}
		}
	}
	return res, nil
}

func ParseTopics(args *Type, topics []block.Hash) ([]interface{}, error) {
	if args.kind != KindTuple {
		return nil, fmt.Errorf("expected a tuple type")
	}
	if len(args.TupleElems()) != len(topics) {
		return nil, fmt.Errorf("bad length")
	}
	elems := make([]interface{}, 0)
	for idx, arg := range args.TupleElems() {
		elem, err := ParseTopic(arg.Elem, topics[idx])
		if err != nil {
			return nil, err
		}
		elems = append(elems, elem)
	}
	return elems, nil
}

var topicTrue, topicFalse block.Hash

func ParseTopic(t *Type, topic block.Hash) (interface{}, error) {
	switch t.kind {
	case KindBool:
		if bytes.Equal(topic[:], topicTrue[:]) {
			return true, nil
		} else if bytes.Equal(topic[:], topicFalse[:]) {
			return false, nil
		}
		return true, fmt.Errorf("is not a boolean")
	case KindInt, KindUInt:
		return decodeInteger(t, topic[:]), nil
	case KindAddress:
		return decodeAddress(topic[:])
	case KindFixedBytes:
		return decodeFixedBytes(t, topic[:])
	default:
		return nil, fmt.Errorf("topic parsing for type '%s' not supported", t.String())
	}
}
