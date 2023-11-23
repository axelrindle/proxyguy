package modules

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/axelrindle/proxyguy/config"
)

type conf = map[string]any

func getConfigFile() string {
	home := os.Getenv("HOME")
	return path.Join(home, ".docker", "config.json")
}

func loadConfig() conf {
	configFile := getConfigFile()
	store := &conf{}

	file, err := os.OpenFile(configFile, os.O_SYNC|os.O_RDONLY, 0)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		os.Stderr.Write([]byte("Failed loading Docker config from " + configFile + "!"))
	}
	defer file.Close()

	json.NewDecoder(file).Decode(store)

	return *store
}

func writeConfig(store *conf) error {
	configFile := getConfigFile()

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

var TemplateDocker = &Module{
	Name:     "Docker",
	Template: "",

	IsEnabled: func(cfg config.StructureModules) bool {
		return cfg.Docker
	},

	Preprocess: func(data Exports) Exports {
		dockerConfig := loadConfig()

		httpString := fmt.Sprintf("http://%s:%s", data.Host, data.Port)
		dockerConfig["proxies"] = map[string]any{
			"default": map[string]string{
				"httpProxy":  httpString,
				"httpsProxy": httpString,
				"noProxy":    data.NoProxy,
			},
		}

		writeConfig(&dockerConfig)

		return data
	},

	OnNoProxy: func() {
		dockerConfig := loadConfig()

		delete(dockerConfig, "proxies")

		writeConfig(&dockerConfig)
	},
}
