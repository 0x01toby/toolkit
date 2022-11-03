package jsonrpc

import "github.com/silenceper/pool"

type PoolCfg struct {
	pool.Config
	headers map[string]string
}
