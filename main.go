package main

import (
	"flag"
	"net/url"
	"os"
	"strings"

	"github.com/axelrindle/proxyguy/config"
	"github.com/axelrindle/proxyguy/modules"
	"github.com/axelrindle/proxyguy/pac"
	"github.com/axelrindle/proxyguy/server"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// flags
type Options struct {
	configFile  string
	printConfig bool
	startServer bool
	verbosity   string
}

func main() {
	os.Unsetenv("http_proxy")
	os.Unsetenv("https_proxy")
	os.Unsetenv("no_proxy")
	os.Unsetenv("HTTP_PROXY")
	os.Unsetenv("HTTPS_PROXY")
	os.Unsetenv("NO_PROXY")

	opts := &Options{}
	flag.StringVar(&opts.configFile, "config", "/etc/proxyguy/config.yaml", "Specify an alternative config file to use.")
	flag.BoolVar(&opts.printConfig, "print-config", false, "Print the effective config and exit.")
	flag.BoolVar(&opts.startServer, "server", false, "Whether to start the integrated proxy server.")
	flag.StringVar(&opts.verbosity, "verbosity", "", "Specify verbosity level.")
	flag.Parse()

	logger := logrus.New()
	logger.SetOutput(os.Stderr)
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logrus.TextFormatter{
		DisableSorting:         true,
		DisableLevelTruncation: true,
		SortingFunc:            func(s []string) {},
	})

	if opts.verbosity != "" {
		switch strings.ToLower(opts.verbosity) {
		case "debug":
			logger.SetLevel(logrus.DebugLevel)
		case "trace":
			logger.SetLevel(logrus.TraceLevel)
		default:
			logger.Warnf("Invalid verbosity level \"%s\"!", opts.verbosity)
		}
	}

	config := &config.Config{Logger: logger, File: &opts.configFile}
	cfg := config.Load()

	if opts.printConfig {
		yaml.NewEncoder(os.Stderr).Encode(cfg)
		return
	}

	if opts.startServer {
		server := &server.Server{Logger: logger, Config: cfg}
		server.Start()
	} else {
		u, err := FindProxy(logger, cfg)
		if err != nil {
			logger.Fatal(err)
		}

		data := &modules.Exports{Host: u.Hostname(), Port: u.Port(), NoProxy: cfg.Proxy.NoProxy}

		mdls := []modules.Module{
			*modules.TemplateMain,
			*modules.TemplateMaven,
		}

		for _, mdl := range mdls {
			if !mdl.IsEnabled(cfg.Modules) {
				continue
			}

			var theData modules.Exports
			if mdl.Preprocess != nil {
				theData = mdl.Preprocess(*data)
			} else {
				theData = *data
			}

			if !modules.Process(mdl, theData) {
				logger.Errorf("Failed parsing template \"%s\"!", mdl.Name)
			}
		}
	}
}

func FindProxy(logger *logrus.Logger, cfg *config.Structure) (*url.URL, error) {
	p := &pac.Pac{Logger: logger, Config: cfg}

	if !p.CheckConnectivity() {
		logger.Debugln("Proxy is inactive. Environment will be left unchanged.")
		os.Exit(0)
	}

	err := p.LoadPacScript()
	if err != nil {
		logger.Fatalln(err)
	}

	parts := p.DetermineProxies(nil)
	if len(parts) == 0 {
		logger.Fatalln("Found no available proxy endpoint!")
	}

	u := pac.TrimProxy(parts[0]) // TODO: Make sure not DIRECT
	logger.Debugln("Using \"" + u + "\" as proxy endpoint.")

	return url.Parse("http://" + u)
}
