package polling

import (
	"context"
	"github.com/taorzhang/toolkit/client"
	"github.com/taorzhang/toolkit/logs"
	"time"
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

func (p *Pipeline) isQuit() bool {
	select {
	case <-p.cancel:
		return true
	default:
		return false
	}
}

func (p *Pipeline) NextBlockHeights(start uint64) *NextPollingAction {
	var iStart, iEnd = start, start + uint64(p.config.Step)
	if p.config.Mode.IsChase() {
		if iStart >= p.config.End {
			return NewNextPollingAction(iStart, iEnd, Done)
		}
		if iEnd >= p.config.End {
			return NewNextPollingAction(iStart, p.config.End, ContinuePolling)
		}
		return NewNextPollingAction(iStart, iEnd, ContinuePolling)
	}
	nodeHeight, err := p.client.BlockNumber(context.Background())
	if err != nil {
		return NewNextPollingAction(iStart, iEnd, WaitingBlocks)
	}
	if iStart >= nodeHeight {
		return NewNextPollingAction(iStart, iEnd, WaitingBlocks)
	}
	if iEnd >= nodeHeight {
		return NewNextPollingAction(iStart, nodeHeight, ContinuePolling)
	}
	return NewNextPollingAction(iStart, iEnd, ContinuePolling)
}

func (p *Pipeline) Run(ctx context.Context) {
	var blockHeight = p.config.Start
	limitCh := make(chan struct{}, p.config.Concurrency)
	defer close(limitCh)
	quitPipeline := make(chan bool)
	defer close(quitPipeline)
	go func() {
		for {
			select {
			case <-ctx.Done():
				// 接受到ctx的cancel信号，准备退出
				log.Warn(ctx, "pipeline is ready to quit", "ctx_err", ctx.Err())
				p.items <- NewCancelItem(context.Background())
				close(p.cancel)
			case item := <-p.items:
				// pipeline有数据，则消费数据
				if item.cancel {
					quitPipeline <- true
					break
				}
				go doItem(item)
			case limitCh <- struct{}{}:
				// 生产数据，将数据写入到pipeline中
				if p.isQuit() {
					<-time.After(time.Second)
					break
				}
				nextAction := p.NextBlockHeights(blockHeight)
				if nextAction.IsWaiting() {
					// 等待出块
					<-time.After(time.Second)
					break
				}
				if nextAction.IsDone() {
					// 任务结束
					p.items <- NewCancelItem(context.Background())
					close(p.cancel)
					break
				}
				go func(m, n uint64) {
					defer func() {
						<-limitCh
					}()
					// pulling data and send to pipeline
					pollingItem(m, n)
				}(nextAction.StartHeight(), nextAction.EndHeight())
			}
		}
	}()
	<-quitPipeline
	log.Warn(ctx, "pipeline quit success")
}

func doItem(item *Item) {}

func pollingItem(iStart, iEnd uint64) {}
