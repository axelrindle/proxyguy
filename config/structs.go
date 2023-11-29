package config

type StructureServer struct {
	Address string `yaml:"address" env:"ADDRESS" env-default:"0.0.0.0"`
	Port    uint   `yaml:"port" env:"PORT" env-default:"1337"`
}

type StructureProxy struct {
	Override  string `yaml:"override" env:"OVERRIDE"`
	NoProxy   string `yaml:"ignore" env:"IGNORE" env-default:"localhost,127.0.0.1"`
	Determine string `yaml:"determine-url" env:"DETERMINE" env-default:"https://ubuntu.com"`
}

type StructureModules struct {
	Main   bool `yaml:"main" env:"MAIN" env-default:"true"`
	Maven  bool `yaml:"maven" env:"MAVEN"`
	Gradle bool `yaml:"gradle" env:"GRADLE"`
	Docker bool `yaml:"docker" env:"DOCKER"`
}

type Structure struct {
	PacUrl     string           `yaml:"pac" env:"PAC"`
	Timeout    uint             `yaml:"timeout" env:"TIMEOUT" env-default:"1000"`
	StatusInfo bool             `yaml:"status-info" env:"STATUS_INFO"`
	Proxy      StructureProxy   `yaml:"proxy" env-prefix:"PROXY_"`
	Server     StructureServer  `yaml:"server" env-prefix:"SERVER_"`
	Modules    StructureModules `yaml:"modules" env-prefix:"MODULES_"`
}
