package modules

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/axelrindle/proxyguy/config"
	"github.com/sirupsen/logrus"
)

type conf = map[string]any

type TemplateDocker struct {
	Logger *logrus.Logger
}

func (t TemplateDocker) GetName() string {
	return "Docker"
}

func (t TemplateDocker) GetTemplate() string {
	return ""
}

func (t TemplateDocker) IsEnabled(cfg config.StructureModules) bool {
	return cfg.Docker
}

func (t TemplateDocker) Preprocess(data *Exports) {
	dockerConfig := t.loadConfig()

	httpString := fmt.Sprintf("http://%s:%s", data.Host, data.Port)
	dockerConfig["proxies"] = map[string]any{
		"default": map[string]string{
			"httpProxy":  httpString,
			"httpsProxy": httpString,
			"noProxy":    data.NoProxy,
		},
	}

	t.writeConfig(&dockerConfig)
}

func (t TemplateDocker) OnNoProxy() {
	dockerConfig := t.loadConfig()

	delete(dockerConfig, "proxies")

	t.writeConfig(&dockerConfig)
}

func (t TemplateDocker) getConfigFile() string {
	home := os.Getenv("HOME")
	return path.Join(home, ".docker", "config.json")
}

func (t TemplateDocker) loadConfig() conf {
	configFile := t.getConfigFile()
	store := &conf{}

	file, err := os.OpenFile(configFile, os.O_SYNC|os.O_RDONLY, 0)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		os.Stderr.Write([]byte("Failed loading Docker config from " + configFile + "!"))
	}
	defer file.Close()

	json.NewDecoder(file).Decode(store)

	return *store
}

func (t TemplateDocker) writeConfig(store *conf) error {
	configFile := t.getConfigFile()

	b, err := json.MarshalIndent(store, "", "    ")
	if err != nil {
		return err
	}

	file, err := os.OpenFile(configFile, os.O_SYNC|os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0)
	if err != nil {
		os.Stderr.Write([]byte("Failed loading Docker config from " + configFile + "!"))
	}
	defer file.Close()

	file.Write(b)

	return nil
}
