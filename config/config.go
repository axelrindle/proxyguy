package config

import (
	"os"

	"github.com/caarlos0/env/v10"
	"github.com/creasty/defaults"
	"gopkg.in/yaml.v3"
)

type config struct {
	File string

	data *Structure
}

var instance *config

func New(file string) *config {
	if instance == nil {
		instance = &config{File: file}
		return instance
	}

	return nil
}

func Data() Structure {
	return *instance.data
}

// Loads the configuration from various sources.
// The following precedence takes effect:
//
//   - 1. Environment variables
//   - 2. File values
//   - 3. Defaults
func (c *config) Load() error {
	var cfg Structure

	// load defaults
	err := defaults.Set(&cfg)
	if err != nil {
		return err
	}

	// load from file
	handle, err := os.OpenFile(c.File, os.O_RDONLY|os.O_SYNC, 0)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		defer handle.Close()

		err = yaml.NewDecoder(handle).Decode(&cfg)
		if err != nil {
			return err
		}
	}

	// load env variables
	err = env.Parse(&cfg)
	if err != nil {
		return err
	}

	c.data = &cfg

	return nil
}
