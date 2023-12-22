package cli

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/axelrindle/proxyguy/cli/commands"
	"github.com/axelrindle/proxyguy/config"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var (
	version   = "dev"
	buildTime = time.Now().Local().Format(time.RFC822)

	logger = logrus.New()
	cfg    = config.New("/etc/proxyguy/config.yaml")
)

var before cli.BeforeFunc = func(ctx *cli.Context) error {
	err := cfg.Load()
	if err != nil {
		return err
	}

	logger.Debugf("Loaded config from %s", cfg.File)

	return nil
}

var app = &cli.App{
	Name:    "proxyguy",
	Usage:   "Dynamic proxy generator.",
	Version: version,

	EnableBashCompletion: true,

	Before: before,

	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "verbosity",
			Usage: "Specify verbosity level.",
			Action: func(ctx *cli.Context, s string) error {
				switch strings.ToLower(s) {
				case "debug":
					logger.SetLevel(logrus.DebugLevel)
				case "trace":
					logger.SetLevel(logrus.TraceLevel)
				default:
					return errors.New("invalid verbosity level \"%s\"")
				}

				return nil
			},
		},
		&cli.PathFlag{
			Name:  "config",
			Usage: "Specify an alternative config file to use.",
			Value: cfg.File,
			Action: func(ctx *cli.Context, p cli.Path) error {
				cfg.File = p

				return nil
			},
		},
	},

	DefaultCommand: "run",
	Commands: []*cli.Command{
		commands.CommandPrintConfig,
		commands.CommandServer,
		commands.CommandRun,
	},
}

func Run() {
	cli.VersionPrinter = func(ctx *cli.Context) {
		println("proxyguy version " + ctx.App.Version)
		println("build time was " + buildTime)
	}

	err := app.Run(os.Args)
	if err != nil {
		logger.Fatal(err)
	}
}
