package polling

import (
	"strings"
	"time"
)

type Mode string

func (m Mode) String() string {
	return string(m)
}

func (m Mode) IsChase() bool {
	return strings.EqualFold(m.String(), ChaseMode.String())
}

var (
	LinearMode = Mode("Liner")
	ChaseMode  = Mode("Chase")
)

type Config struct {
	ItemLen      uint
	Start        uint64
	End          uint64
	Mode         Mode
	SkipInternal bool
	// Step 每次取块数量
	Step uint
	// Concurrency 并发
	Concurrency         uint
	QuitWaitingDuration time.Duration
}

func NewLinerConfig() *Config {
	return &Config{ItemLen: 100, Step: 8, Concurrency: 1, Mode: LinearMode, QuitWaitingDuration: 10 * time.Second}
}

func NewChaseConfig() *Config {
	return &Config{ItemLen: 100, Step: 16, Concurrency: 3, Mode: ChaseMode, QuitWaitingDuration: 30 * time.Second}
}

type CfgOpt func(c *Config) error
