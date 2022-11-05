package jsonrpc

import "testing"

func Test_explode(t *testing.T) {
	var arr []int
	for i := 0; i < 100; i++ {
		arr = append(arr, i)
	}

	size := explodeBySize(arr, 21)
	for idx := range size {
		t.Log("idx", idx, "items", size[idx])
	}
}
