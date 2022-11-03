package jsonrpc

import (
	"github.com/ethereum/go-ethereum/rpc"
	"time"
)

var (
	InitCap     = 5
	MaxIdle     = 20
	MaxCap      = 100
	IdleTimeout = 5 * time.Second
)

type PoolCfgOpt func(c *PoolCfg)

func GetDefaultOpts(endpoint string) []PoolCfgOpt {
	return []PoolCfgOpt{
		WithRpcFactory(endpoint),
		WithRpcClose(),
		WithInitCap(InitCap),
		WithMaxIdle(MaxIdle),
		WithMaxCap(MaxCap),
		WithIdleTimeout(IdleTimeout),
	}
}

func GetEthCfgOpts(endpoint string, initCap, maxCap, maxIdleCap int, maxIdle time.Duration) []PoolCfgOpt {
	return []PoolCfgOpt{
		WithRpcFactory(endpoint),
		WithRpcClose(),
		WithInitCap(initCap),
		WithMaxIdle(maxIdleCap),
		WithMaxCap(maxCap),
		WithIdleTimeout(maxIdle),
	}
}

// WithRpcHeaders headers
func WithRpcHeaders(headers map[string]string) PoolCfgOpt {
	return func(c *PoolCfg) {
		c.headers = headers
	}
}

// WithRpcFactory factory
func WithRpcFactory(endpoint string) PoolCfgOpt {
	return func(c *PoolCfg) {
		c.Factory = func() (interface{}, error) {
			return rpc.Dial(endpoint)
		}
	}
}

// WithRpcClose rpc close时
func WithRpcClose() PoolCfgOpt {
	return func(c *PoolCfg) {
		c.Close = func(v interface{}) error {
			v.(*rpc.Client).Close()
			return nil
		}
	}
}

func WithInitCap(cap int) PoolCfgOpt {
	return func(c *PoolCfg) {
		c.InitialCap = cap
	}
}

func WithMaxCap(cap int) PoolCfgOpt {
	return func(c *PoolCfg) {
		c.MaxCap = cap
	}
}

func WithMaxIdle(cap int) PoolCfgOpt {
	return func(c *PoolCfg) {
		c.MaxIdle = cap
	}
}

func WithIdleTimeout(t time.Duration) PoolCfgOpt {
	return func(c *PoolCfg) {
		c.IdleTimeout = t
	}
}
