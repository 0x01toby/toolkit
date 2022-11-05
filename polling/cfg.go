package polling

type Mode string

var (
	LinearMode = Mode("Liner")
	ChaseMode  = Mode("Chase")
)

type Config struct {
	ItemLen      uint
	Start        uint
	End          uint
	Mode         Mode
	SkipInternal bool
	// Step 每次取块数量
	Step uint
	// Concurrency 并发
	Concurrency uint
}

func NewConfig() *Config {
	return &Config{ItemLen: 100, Step: 8, Concurrency: 5, Mode: LinearMode}
}

type CfgOpt func(c *Config) error
