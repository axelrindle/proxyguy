package config

import (
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type StructureServer struct {
	Address *string `yaml:"address"`
	Port    uint    `yaml:"port"`
}

type StructureProxy struct {
	Override  string `yaml:"override"`
	NoProxy   string `yaml:"ignore"`
	Determine string `yaml:"determine-url"`
}

type Structure struct {
	PacUrl  string          `yaml:"pac"`
	Timeout uint            `yaml:"timeout"`
	Proxy   StructureProxy  `yaml:"proxy"`
	Server  StructureServer `yaml:"server"`
}

type Config struct {
	Logger *logrus.Logger
	File   *string
}

func (c Config) buildDefault() *Structure {
	address := "0.0.0.0"

	return &Structure{
		Timeout: 1000,
		Proxy: StructureProxy{
			NoProxy:   "localhost,127.0.0.1",
			Determine: "https://ubuntu.com",
		},
		Server: StructureServer{
			Address: &address,
			Port:    1337,
		},
	}
}

func (c Config) Load() *Structure {
	f, err := os.Open(*c.File)
	if err != nil {
		c.Logger.Fatalln(err)
	}
	defer f.Close()

	cfg := c.buildDefault()
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		c.Logger.Fatalln(err)
	}

	return cfg
}
