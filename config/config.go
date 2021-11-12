package config

import (
	"github.com/kkyr/fig"
	"log"
	"time"
)

// Config structure for settings of application
type Config struct {
	App struct {
		URL            string        `fig:"URL" default:"http://home.mcom.com/home/welcome.html"` // address of target source
		TimeoutRequest time.Duration `fig:"timeoutRequest" default:"10"`                          // request timeout in seconds
		TimeoutApp     time.Duration `fig:"timeoutApp" default:"180"`                             // application timeout in seconds
		MaxDepth       uint          `fig:"maxDepth" default:"3"`                                 // max depth for links
		MaxResults     uint          `fig:"maxResults" default:"500"`                             // max result of links
		MaxErrors      uint          `fig:"maxErrors" default:"500"`                              // max errors of request results
	} `fig:"app"`
}

// Init function for initialize Config structure
func Init() (*Config, error) {
	var cfg = Config{}
	err := fig.Load(&cfg, fig.Dirs("config/", "../config/", "../../config/", "../../../config/"), fig.File("config.yaml"))
	if err != nil {
		log.Fatalf("can't load configuration file: %s", err)
		return nil, err
	}

	return &cfg, err
}
