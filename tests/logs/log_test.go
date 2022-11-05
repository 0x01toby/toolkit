package logs

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/taorzhang/toolkit/logs"
	"github.com/taorzhang/toolkit/logs/tracing"
	"testing"
	"time"
)

var (
	log = logs.NewLogger("module", "logs")
)

func TestLog(t *testing.T) {
	log.Info(context.Background(), "hello world", "a", 1, "b", 2)
}

func nodeRpc(ctx context.Context) {
	log.Info(ctx, "模拟开始请求node rpc 数据")
	_, span := tracing.StartSpan(ctx, "retrieve_node", "txs")
	defer span.End()
	select {
	case <-time.After(3 * time.Second):
		log.Info(ctx, "模拟结束请求node rpc 数据")
		return
	}
}

func TestJaeger(t *testing.T) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	cfg := tracing.NewTraceCfg("http://127.0.0.1:14268/api/traces", "bcdb", "eth")
	err := tracing.InitOpenTrace(ctx, cfg, tracing.WithSampler(1))
	assert.NoError(t, err)
	ctx, span := tracing.StartSpan(context.Background(), "retrieve_node", "blocks2")
	nodeRpc(ctx)
	span.End()
	log.Info(ctx, "测试jaeger, retrieve node done.")
	cancelFunc()
	// 退出时，需要更多的时间上报数据
	<-time.After(3 * time.Second)
}
