package config

import (
	"log"
	"testing"
)

func TestInit(t *testing.T) {
	plug := "config.yaml"
	cfg, err := InitConfig(&plug)
	if cfg == nil {
		log.Fatal("error: can't load configuration")
	}

	if err != nil {
		log.Fatalf("error: can't load configuration: %s", err)
	}
}
