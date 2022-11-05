package polling

import (
	"context"
	"github.com/taorzhang/toolkit/client"
	"github.com/taorzhang/toolkit/logs"
)

var (
	log = logs.NewLogger("module", "pipeline")
)

// Pipeline blockchain polling pipeline
type Pipeline struct {
	items  chan *Item
	cancel chan bool
	config *Config
	client client.Provider
}

func NewPipeline(opts ...CfgOpt) (*Pipeline, error) {
	c := NewConfig()
	for idx := range opts {
		if err := opts[idx](c); err != nil {
			return nil, err
		}
	}
	return &Pipeline{items: make(chan *Item, c.ItemLen), cancel: make(chan bool), config: c}, nil
}

func (p *Pipeline) Run(ctx context.Context) {
	go func() {
		select {
		// 准备退出
		case <-ctx.Done():
			log.Warn(ctx, "pipeline is ready to quit", "ctx_err", ctx.Err())
			p.cancel <- true
			return
		}
	}()

	limitCh := make(chan struct{}, p.config.Concurrency)
	defer close(limitCh)
	for {
		select {
		// pipeline有数据，则消费数据
		case item := <-p.items:
			go doItem(item)
		case limitCh <- struct{}{}:
			go func() {
				defer func() {
					<-limitCh
				}()
				// pulling data and send to pipeline
			}()
		default:

		}
	}

}

func doItem(item *Item) {}
