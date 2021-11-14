package config

import (
	"flag"
	"github.com/kkyr/fig"
	"log"
	"time"
)

const (
	usageConfig = "use this flag for set path to configuration file"
)

// Config structure for settings of application
type Config struct {
	App struct {
		URL            string        `fig:"URL" default:"http://home.mcom.com/home/welcome.html"` // address of target source
		TimeoutRequest time.Duration `fig:"timeoutRequest" default:"10"`                          // request timeout in seconds
		TimeoutApp     time.Duration `fig:"timeoutApp" default:"180"`                             // application timeout in seconds
		MaxDepth       uint64        `fig:"maxDepth" default:"3"`                                 // max depth for links
		MaxResults     uint          `fig:"maxResults" default:"500"`                             // max result of links
		MaxErrors      uint          `fig:"maxErrors" default:"500"`                              // max errors of request results
		DeltaDepth     uint64        `fig:"deltaDepth" default:"2"`                               // delta for increment depth
	} `fig:"app"`
}

// Init function for initialize Config structure
func Init() (*Config, error) {
	useConfig := flag.String("path", "config/config.yaml", usageConfig)
	flag.Parse()

	var cfg = Config{}
	err := fig.Load(&cfg, fig.File(*useConfig))
	if err != nil {
		log.Fatalf("can't load configuration file: %s", err)
		return nil, err
	}

	return &cfg, err
}
