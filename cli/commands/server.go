package commands

import (
	"github.com/axelrindle/proxyguy/server"
	"github.com/urfave/cli/v2"
)

var CommandServer = &cli.Command{
	Name:  "server",
	Usage: "Start the integrated proxy server.",
	Action: func(ctx *cli.Context) error {
		server := &server.Server{}
		server.Start()

		return nil
	},
}
