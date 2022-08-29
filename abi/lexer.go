package abi

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
)

type token struct {
	typ     tokenType
	literal string
}

func (t tokenType) String() string {
	names := [...]string{
		"eof",
		"string",
		"number",
		"tuple",
		"(",
		")",
		"[",
		"]",
		",",
		"indexed",
		"<invalid>",
	}
	return names[t]
}

func expectedToken(t tokenType) error {
	return fmt.Errorf("expected token %s", t.String())
}

func notExpectedToken(t tokenType) error {
	return fmt.Errorf("token '%s' not expected", t.String())
}

// lexer 词法分析
type lexer struct {
	input        string
	current      token
	peek         token
	position     int
	readPosition int
	ch           byte
}

func newLexer(input string) *lexer {
	l := &lexer{input: input}
	l.readChar()
	return l
}

func (l *lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition++
}

func (l *lexer) nextToken() token {
	l.current = l.peek
	l.peek = l.nextTokenImpl()
	return l.current
}

func (l *lexer) nextTokenImpl() token {
	var tok token

	// skip whitespace
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}

	switch l.ch {
	case ',':
		tok.typ = commaToken
	case '(':
		tok.typ = lparenToken
	case ')':
		tok.typ = rparenToken
	case '[':
		tok.typ = lbracketToken
	case ']':
		tok.typ = rbracketToken
	case 0:
		tok.typ = eofToken
	default:
		if isLetter(l.ch) {
			tok.literal = l.readIdentifier()
			if tok.literal == "tuple" {
				tok.typ = tupleToken
			} else if tok.literal == "indexed" {
				tok.typ = indexedToken
			} else {
				tok.typ = strToken
			}
			return tok
		} else if isDigit(l.ch) {
			return token{numberToken, l.readNumber()}
		} else {
			tok.typ = invalidToken
		}
	}

	l.readChar()
	return tok
}

func (l *lexer) readIdentifier() string {
	pos := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}

	return l.input[pos:l.position]
}

func (l *lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func readType(l *lexer) (*Type, error) {
	var tt *Type

	tok := l.nextToken()

	isTuple := false
	if tok.typ == tupleToken {
		if l.nextToken().typ != lparenToken {
			return nil, expectedToken(lparenToken)
		}
		isTuple = true
	} else if tok.typ == lparenToken {
		isTuple = true
	}
	if isTuple {
		var next token
		elems := []*TupleElem{}
		for {

			name := ""
			indexed := false

			elem, err := readType(l)
			if err != nil {
				if l.current.typ == rparenToken && len(elems) == 0 {
					// empty tuple 'tuple()'
					break
				}
				return nil, fmt.Errorf("failed to decode type: %v", err)
			}

			switch l.peek.typ {
			case strToken:
				l.nextToken()
				name = l.current.literal

			case indexedToken:
				l.nextToken()
				indexed = true
				if l.peek.typ == strToken {
					l.nextToken()
					name = l.current.literal
				}
			}

			elems = append(elems, &TupleElem{
				Name:    name,
				Elem:    elem,
				Indexed: indexed,
			})

			next = l.nextToken()
			if next.typ == commaToken {
				continue
			} else if next.typ == rparenToken {
				break
			} else {
				return nil, notExpectedToken(next.typ)
			}
		}
		tt = &Type{kind: KindTuple, tuple: elems, t: tupleT}

	} else if tok.typ != strToken {
		return nil, expectedToken(strToken)

	} else {
		// Check normal types
		elem, err := decodeSimpleType(tok.literal)
		if err != nil {
			return nil, err
		}
		tt = elem
	}

	// check for arrays at the end of the type
	for {
		if l.peek.typ != lbracketToken {
			break
		}

		l.nextToken()
		n := l.nextToken()

		var tAux *Type
		if n.typ == rbracketToken {
			tAux = &Type{kind: KindSlice, elem: tt, t: reflect.SliceOf(tt.t)}

		} else if n.typ == numberToken {
			size, err := strconv.ParseUint(n.literal, 10, 32)
			if err != nil {
				return nil, fmt.Errorf("failed to read array size '%s': %v", n.literal, err)
			}

			tAux = &Type{kind: KindArray, elem: tt, size: int(size), t: reflect.ArrayOf(int(size), tt.t)}
			if l.nextToken().typ != rbracketToken {
				return nil, expectedToken(rbracketToken)
			}
		} else {
			return nil, notExpectedToken(n.typ)
		}

		tt = tAux
	}
	return tt, nil
}

var typeRegexp = regexp.MustCompile("^([[:alpha:]]+)([[:digit:]]*)$")

func decodeSimpleType(str string) (*Type, error) {
	match := typeRegexp.FindStringSubmatch(str)
	if len(match) == 0 {
		return nil, fmt.Errorf("type format is incorrect. Expected 'type''bytes' but found '%s'", str)
	}
	match = match[1:]

	var err error
	t := match[0]

	bytes := 0
	ok := false

	if bytesStr := match[1]; bytesStr != "" {
		bytes, err = strconv.Atoi(bytesStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse bytes '%s': %v", bytesStr, err)
		}
		ok = true
	}

	// int and uint without bytes default to 256, 'bytes' may
	// have or not, the rest dont have bytes
	if t == "int" || t == "uint" {
		if !ok {
			bytes = 256
		}
	} else if t != "bytes" && ok {
		return nil, fmt.Errorf("type %s does not expect bytes", t)
	}

	switch t {
	case "uint":
		var k reflect.Type
		switch bytes {
		case 8:
			k = uint8T
		case 16:
			k = uint16T
		case 32:
			k = uint32T
		case 64:
			k = uint64T
		default:
			if bytes%8 != 0 {
				panic(fmt.Errorf("number of bytes has to be M mod 8"))
			}
			k = bigIntT
		}
		return &Type{kind: KindUInt, size: int(bytes), t: k}, nil

	case "int":
		var k reflect.Type
		switch bytes {
		case 8:
			k = int8T
		case 16:
			k = int16T
		case 32:
			k = int32T
		case 64:
			k = int64T
		default:
			if bytes%8 != 0 {
				panic(fmt.Errorf("number of bytes has to be M mod 8"))
			}
			k = bigIntT
		}
		return &Type{kind: KindInt, size: int(bytes), t: k}, nil

	case "byte":
		bytes = 1
		fallthrough

	case "bytes":
		if bytes == 0 {
			return &Type{kind: KindBytes, t: dynamicBytesT}, nil
		}
		return &Type{kind: KindFixedBytes, size: int(bytes), t: reflect.ArrayOf(int(bytes), reflect.TypeOf(byte(0)))}, nil

	case "string":
		return &Type{kind: KindString, t: stringT}, nil

	case "bool":
		return &Type{kind: KindBool, t: boolT}, nil

	case "address":
		return &Type{kind: KindAddress, t: addressT, size: 20}, nil

	case "function":
		return &Type{kind: KindFunction, size: 24, t: functionT}, nil

	default:
		return nil, fmt.Errorf("unknown type '%s'", t)
	}
}

type tokenType int

const (
	eofToken tokenType = iota
	strToken
	numberToken
	tupleToken
	lparenToken
	rparenToken
	lbracketToken
	rbracketToken
	commaToken
	indexedToken
	invalidToken
)
