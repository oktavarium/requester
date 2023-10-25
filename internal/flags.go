package requester

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env"
)

type config struct {
	FilePath string `env:"FILEPATH"`
}

func loadConfig() (config, error) {
	var flagsConfig config
	flag.StringVar(&flagsConfig.FilePath,
		"f",
		"",
		"path to file with urls")
	flag.Parse()

	if err := env.Parse(&flagsConfig); err != nil {
		return flagsConfig, fmt.Errorf("error occured on parsing env: %w", err)
	}

	if len(flag.Args()) > 0 {
		return flagsConfig, fmt.Errorf("unrecognised flags found")
	}

	if len(flagsConfig.FilePath) == 0 {
		return flagsConfig, fmt.Errorf("file path can't be empty")
	}

	return flagsConfig, nil
}
