package jsonrpc

import (
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/silenceper/pool"
	"time"
)

type PoolCfgOpt func(c *pool.Config)

func GetEthCfgOpts(endpoint string, initCap, maxCap, maxIdleCap int, maxIdle time.Duration) []PoolCfgOpt {
	return []PoolCfgOpt{
		WithEthRpcFactory(endpoint),
		WithEthRpcClose(),
		WithInitCap(initCap),
		WithMaxIdle(maxIdleCap),
		WithMaxCap(maxCap),
		WithIdleTimeout(maxIdle),
	}
}

func WithEthRpcFactory(endpoint string) PoolCfgOpt {
	return func(c *pool.Config) {
		c.Factory = func() (interface{}, error) {
			return rpc.Dial(endpoint)
		}
	}
}

func WithEthRpcClose() PoolCfgOpt {
	return func(c *pool.Config) {
		c.Close = func(v interface{}) error {
			v.(*rpc.Client).Close()
			return nil
		}
	}
}

func WithInitCap(cap int) PoolCfgOpt {
	return func(c *pool.Config) {
		c.InitialCap = cap
	}
}

func WithMaxCap(cap int) PoolCfgOpt {
	return func(c *pool.Config) {
		c.MaxCap = cap
	}
}

func WithMaxIdle(cap int) PoolCfgOpt {
	return func(c *pool.Config) {
		c.MaxIdle = cap
	}
}

func WithIdleTimeout(t time.Duration) PoolCfgOpt {
	return func(c *pool.Config) {
		c.IdleTimeout = t
	}
}
