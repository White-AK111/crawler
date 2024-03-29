package config

import (
	"log"
	"sync/atomic"
	"time"

	"github.com/kkyr/fig"
)

// Config structure for settings of application
type Config struct {
	App struct {
		URL            string        `fig:"URL" default:"https://golang.org"` // address of target source
		TimeoutRequest time.Duration `fig:"timeoutRequest" default:"10"`      // request timeout in seconds
		TimeoutApp     time.Duration `fig:"timeoutApp" default:"180"`         // application timeout in seconds
		MaxDepth       int64         `fig:"maxDepth" default:"3"`             // max depth for links
		MaxResults     uint          `fig:"maxResults" default:"500"`         // max result of links
		MaxErrors      uint          `fig:"maxErrors" default:"500"`          // max errors of request results
		DeltaDepth     int64         `fig:"deltaDepth" default:"2"`           // delta for increment depth
	} `fig:"app"`
}

// InitConfig function for initialize Config structure
func InitConfig(useConfig *string) (*Config, error) {
	var cfg = Config{}
	err := fig.Load(&cfg, fig.File(*useConfig))
	if err != nil {
		err = fig.Load(&cfg, fig.File("config.yaml"))
		if err != nil {
			log.Fatalf("can't load configuration file: %s", err)
			return nil, err
		}
	}

	return &cfg, err
}

// ChangeMaxDepth increment depth by atomic
func (c *Config) ChangeMaxDepth(delta int64) {
	atomic.AddInt64(&c.App.MaxDepth, delta)
}
