package jsonrpc

import (
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/silenceper/pool"
)

type Pool struct {
	pool.Pool
	headers map[string]string
}

func NewPool(headers map[string]string, opts ...PoolCfgOpt) (*Pool, error) {
	poolConfig := &pool.Config{}
	for _, opt := range opts {
		opt(poolConfig)
	}
	channelPool, err := pool.NewChannelPool(poolConfig)
	if err != nil {
		return nil, err
	}
	return &Pool{Pool: channelPool, headers: headers}, nil
}

func (p *Pool) GetClient() (*rpc.Client, error) {
	c, err := p.Get()
	if err != nil {
		return nil, err
	}
	client := c.(*rpc.Client)
	if len(p.headers) > 0 {
		for k, v := range p.headers {
			client.SetHeader(k, v)
		}
	}
	return client, err
}

func (p *Pool) PutClient(client *rpc.Client) {
	_ = p.Put(client)
}

func (p *Pool) Run(runnable func(client *rpc.Client) error) error {
	client, err := p.GetClient()
	if err != nil {
		return err
	}
	defer p.PutClient(client)
	return runnable(client)
}
