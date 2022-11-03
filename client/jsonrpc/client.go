package jsonrpc

import (
	"github.com/ethereum/go-ethereum/rpc"
	"sync"
)

var (
	GroupSize = 50
)

type Client struct {
	pool *Pool
}

func NewClient(opts ...PoolCfgOpt) (*Client, error) {
	newPool, err := NewPool(opts...)
	return &Client{pool: newPool}, err
}

func (c *Client) Release() {
	c.pool.Release()
}

// Call 单独call
func (c *Client) Call(method string, out interface{}, args ...interface{}) error {
	return c.pool.Run(func(client *rpc.Client) error {
		return client.Call(out, method, args...)
	})
}

// BatchCall 批量rpc请求，当批量数过多，会进行分组
func (c *Client) BatchCall(elems []rpc.BatchElem, allOk bool) error {
	if len(elems) == 0 {
		return nil
	}
	segments := explodeBySize(elems, int64(GroupSize))
	var wg sync.WaitGroup
	groupErrs := make([]error, len(elems)/GroupSize+1)
	for idx, segment := range segments {
		wg.Add(1)
		go func(idx int, batch []rpc.BatchElem) {
			defer wg.Done()
			client, err := c.pool.GetClient()
			if err != nil {
				return
			}
			defer c.pool.PutClient(client)
			groupErrs[idx] = client.BatchCall(batch)
		}(idx, segment)
	}
	wg.Wait()
	for idx := range groupErrs {
		if groupErrs[idx] != nil {
			return groupErrs[idx]
		}
	}
	if allOk {
		for idx := range elems {
			if elems[idx].Error != nil {
				return elems[idx].Error
			}
		}
	}
	return nil
}

func explodeBySize[T any](arr []T, groupSize int64) [][]T {
	var segments = make([][]T, 0)
	elemNumbers := int64(len(arr))
	if elemNumbers < groupSize {
		segments = append(segments, arr)
		return segments
	}
	groupNumber := elemNumbers / groupSize
	for groupIdx := int64(0); groupIdx <= groupNumber; groupIdx++ {
		start := groupIdx * groupSize
		end := (groupIdx + 1) * groupSize
		if start == elemNumbers {
			break
		}
		if end > elemNumbers {
			end = elemNumbers
		}
		segments = append(segments, arr[start:end])
	}
	return segments
}
