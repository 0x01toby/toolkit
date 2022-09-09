package wallet

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDecimal(t *testing.T) {
	decimal, err := ToWei("0.1")
	assert.NoError(t, err)
	t.Log("decimal", decimal)
}

func TestA(t *testing.T) {
	wei, err := ToWei("1")
	assert.NoError(t, err)
	t.Log("wei", wei)
}

func TestToWei(t *testing.T) {
	testCases := []struct {
		Num    string
		Result string
	}{
		{
			Num:    "1",
			Result: "1000000000000000000",
		},
		{
			Num:    "0.05",
			Result: "50000000000000000",
		},
		{
			Num:    "100",
			Result: "100000000000000000000",
		},
	}
	for _, testCase := range testCases {
		decimal, err := ToWei(testCase.Num)
		assert.NoError(t, err)
		assert.Equal(t, testCase.Result, decimal.String())
	}
}
