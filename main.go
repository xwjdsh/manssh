package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Usage = usage
	app.Version = version
	app.Flags = flags()
	app.Commands = commands()
	app.Run(os.Args)
}

func flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{Name: "file, f", Value: fmt.Sprintf("%s/.ssh/config", getHomeDir()), Destination: &path},
	}
}

func commands() []cli.Command {
	return []cli.Command{
		{
			Name:    "add",
			Usage:   "adds a new ssh alias record",
			Action:  addAction,
			Aliases: []string{"a"},
			Flags: []cli.Flag{
				cli.GenericFlag{Name: "config, c", Value: &kvFlag{}},
			},
		},
		{
			Name:    "list",
			Usage:   "list or search ssh alias records",
			Action:  listAction,
			Aliases: []string{"l"},
		},
		{
			Name:    "update",
			Usage:   "update existing ssh alias record",
			Action:  updateAction,
			Aliases: []string{"u"},
			Flags: []cli.Flag{
				cli.GenericFlag{Name: "config, c", Value: &kvFlag{}},
				cli.StringFlag{Name: "rename, r"},
			},
		},
		{
			Name:    "delete",
			Usage:   "delete existing ssh alias record",
			Action:  deleteAction,
			Aliases: []string{"d"},
		},
		{
			Name:    "backup",
			Usage:   "backup ssh alias config records",
			Action:  backupAction,
			Aliases: []string{"b"},
		},
	}
}
