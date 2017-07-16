package main

import (
	"github.com/urfave/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "asd"
	app.Version = "0.1.0"
	app.Usage = "A Shcript Downloader"
	app.Commands = []cli.Command{
		{
			Name:   "install",
			Usage:  "Installs given modules into the remote server",
			Action: CliInstall,
			Flags: []cli.Flag{
				cli.StringFlag{Name: "remote, r", Usage: "Install into a remote instance via ssh"},
			},
		},
		{
			Name:   "gen",
			Usage:  "Generates env file for given modules for custom variable substitution",
			Action: CliGen,
			Flags: []cli.Flag{
				cli.StringFlag{Name: "file, f", Usage: "Name of generated file"},
			},
		},
		{
			Name:   "update-cache",
			Usage:  "Updates the cache",
			Action: CliUpdateCache,
		},
	}

	app.Run(os.Args)
}
