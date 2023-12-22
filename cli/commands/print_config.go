package commands

import (
	"os"

	"github.com/axelrindle/proxyguy/config"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

var CommandPrintConfig = &cli.Command{
	Name:  "print-config",
	Usage: "Print the effective config and exit.",
	Action: func(ctx *cli.Context) error {
		yaml.NewEncoder(os.Stderr).Encode(config.Data())
		return nil
	},
}
