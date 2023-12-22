package config

type StructureServer struct {
	Address string `yaml:"address" env:"SERVER_ADDRESS" default:"0.0.0.0"`
	Port    uint   `yaml:"port" env:"SERVER_PORT" default:"1337"`
}

type StructureProxy struct {
	Override  string `yaml:"override" env:"PROXY_OVERRIDE"`
	NoProxy   string `yaml:"ignore" env:"PROXY_IGNORE" default:"localhost,127.0.0.1"`
	Determine string `yaml:"determine-url" env:"PROXY_DETERMINE" default:"https://ubuntu.com"`
}

type StructureModules struct {
	Main   bool `yaml:"main" env:"MODULES_MAIN" default:"true"`
	Maven  bool `yaml:"maven" env:"MODULES_MAVEN"`
	Gradle bool `yaml:"gradle" env:"MODULES_GRADLE"`
	Docker bool `yaml:"docker" env:"MODULES_DOCKER"`
}

type Structure struct {
	PacUrl     string `yaml:"pac" env:"PAC"`
	Timeout    uint   `yaml:"timeout" env:"TIMEOUT" default:"1000"`
	StatusInfo bool   `yaml:"status-info" env:"STATUS_INFO"`

	Proxy   StructureProxy   `yaml:"proxy"`
	Server  StructureServer  `yaml:"server"`
	Modules StructureModules `yaml:"modules"`
}
