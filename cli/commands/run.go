package commands

import (
	"errors"
	"net/url"

	"github.com/axelrindle/proxyguy/config"
	"github.com/axelrindle/proxyguy/logger"
	"github.com/axelrindle/proxyguy/modules"
	"github.com/axelrindle/proxyguy/pac"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var ErrNoProxy = errors.New("NO_PROXY")

var (
	log = logger.Logger.WithField("command", "run")
)

func findProxy(cfg config.Structure) (*url.URL, error) {
	p := &pac.Pac{Config: cfg}

	logLevel := logrus.DebugLevel
	if cfg.StatusInfo {
		logLevel = logrus.InfoLevel
	}

	if !p.CheckConnectivity() {
		log.Logln(logLevel, "Proxy is inactive. Environment will be left unchanged.")
		return nil, ErrNoProxy
	}

	err := p.LoadPacScript()
	if err != nil {
		log.Fatalln(err)
	}

	parts := p.DetermineProxies(nil)
	if len(parts) == 0 {
		log.Fatalln("Found no available proxy endpoint!")
	}

	u := pac.TrimProxy(parts[0]) // TODO: Make sure not DIRECT
	log.Logln(logLevel, "Using \""+u+"\" as proxy endpoint.")

	return url.Parse("http://" + u)
}

var CommandRun = &cli.Command{
	Name:  "run",
	Usage: "Generate proxy entries.",
	Action: func(ctx *cli.Context) error {
		mdls := []modules.Module{
			modules.TemplateMain{},
			modules.TemplateMaven{},
			modules.TemplateGradle{},
			modules.TemplateDocker{},
		}

		cfg := config.Data()
		u, err := findProxy(cfg)
		if err == ErrNoProxy {
			for _, mdl := range mdls {
				if !mdl.IsEnabled(cfg.Modules) {
					continue
				}

				mdl.OnNoProxy()
			}

			return nil
		} else if err != nil {
			log.Fatal(err)
		}

		data := &modules.Exports{Host: u.Hostname(), Port: u.Port(), NoProxy: cfg.Proxy.NoProxy}

		for _, mdl := range mdls {
			if !mdl.IsEnabled(cfg.Modules) {
				continue
			}

			moduleData := *data

			mdl.GetLogger().Debug("Running pre-processor")
			mdl.Preprocess(&moduleData)

			mdl.GetLogger().Debug("Running main processor")
			if !modules.Process(mdl, moduleData) {
				log.Errorf("Failed parsing template \"%s\"!", mdl.GetName())
			}
		}

		return nil
	},
}
