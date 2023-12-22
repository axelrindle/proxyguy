package config

import (
	"os"

	"github.com/caarlos0/env"
	"github.com/creasty/defaults"
	"gopkg.in/yaml.v3"
)

type Config struct {
	File string
}

// Loads the configuration from various sources.
// The following precedence takes effect:
//
//   - 1. Environment variables
//   - 2. File values
//   - 3. Defaults
func (c *Config) Load() (*Structure, error) {
	var cfg Structure

	// load defaults
	err := defaults.Set(&cfg)
	if err != nil {
		return nil, err
	}

	// load from file
	handle, err := os.OpenFile(c.File, os.O_RDONLY|os.O_SYNC, 0)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
	} else {
		defer handle.Close()

		err = yaml.NewDecoder(handle).Decode(&cfg)
		if err != nil {
			return nil, err
		}
	}

	// load env variables
	err = env.Parse(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
