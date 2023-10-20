package config

import (
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Logger *logrus.Logger
	File   *string
}

func (c *Config) Load() *Structure {
	var cfg Structure

	err := cleanenv.ReadConfig(*c.File, &cfg)
	if err != nil {
		if os.IsNotExist(err) {
			cleanenv.ReadEnv(&cfg)
		} else {
			c.Logger.Fatalln(err)
		}
	}

	return &cfg
}
