package abi

import (
	"fmt"
	"regexp"
	"strings"
)

// parseEventSignature 解析Event
func parseEventSignature(sig string) (string, *Type, error) {
	sig = strings.TrimPrefix(sig, "event ")
	if !strings.HasSuffix(sig, ")") {
		return "", nil, fmt.Errorf("failed to parse input, expected 'name(types)'")
	}
	idx := strings.Index(sig, "(")
	if idx == -1 {
		return "", nil, fmt.Errorf("failed to parse input, expected 'name(types)'")
	}
	eventName, signature := sig[:idx], sig[idx:]
	signature = "tuple" + signature
	typ, err := NewType(signature)
	if err != nil {
		return "", nil, err
	}
	return eventName, typ, nil
}

var (
	funcRegexpWithReturn    = regexp.MustCompile(`(\w*)\s*\((.*)\)(.*)\s*returns\s*\((.*)\)`)
	funcRegexpWithoutReturn = regexp.MustCompile(`(\w*)\s*\((.*)\)(.*)`)
)

// parseMethodSignature 解析Method
func parseMethodSignature(sig string) (string, *Type, *Type, error) {
	sig = strings.Replace(sig, "\n", " ", -1)
	sig = strings.Replace(sig, "\t", " ", -1)

	sig = strings.TrimPrefix(sig, "function ")
	sig = strings.TrimSpace(sig)

	var funcName, inputArgs, outputArgs string
	if strings.Contains(sig, "returns") {
		matches := funcRegexpWithReturn.FindAllStringSubmatch(sig, -1)
		if len(matches) == 0 {
			return "", nil, nil, fmt.Errorf("no matches found")
		}
		funcName = strings.TrimSpace(matches[0][1])
		inputArgs = strings.TrimSpace(matches[0][2])
		outputArgs = strings.TrimSpace(matches[0][4])
	} else {
		matches := funcRegexpWithoutReturn.FindAllStringSubmatch(sig, -1)
		if len(matches) == 0 {
			return "", nil, nil, fmt.Errorf("no matches found")
		}
		funcName = strings.TrimSpace(matches[0][1])
		inputArgs = strings.TrimSpace(matches[0][2])
	}
	input, err := NewType("tuple(" + inputArgs + ")")
	if err != nil {
		return "", nil, nil, err
	}
	output, err := NewType("tuple(" + outputArgs + ")")
	if err != nil {
		return "", nil, nil, err
	}
	return funcName, input, output, nil
}

func parseErrorSignature(sig string) (string, *Type, error) {
	sig = strings.TrimPrefix(sig, "error ")
	if !strings.HasSuffix(sig, ")") {
		return "", nil, fmt.Errorf("failed to parse input, expected 'name(types)'")
	}
	idx := strings.Index(sig, "(")
	if idx == -1 {
		return "", nil, fmt.Errorf("failed to parse input, expected 'name(types)'")
	}
	funcName, signature := sig[:idx], sig[idx:]
	signature = "tuple" + signature
	typ, err := NewType(signature)
	if err != nil {
		return "", nil, err
	}
	return funcName, typ, nil
}
