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
			Usage:   "add a new ssh alias record",
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
		{
			Name:    "open",
			Usage:   "open new terminal and connecting server only for osx",
			Action:  openAction,
			Aliases: []string{"o"},
		},
		{
			Name:    "run",
			Usage:   "run command on remote server",
			Action:  runAction,
			Aliases: []string{"r"},
			Flags: []cli.Flag{
				cli.BoolFlag{Name: "password, p"},
			},
		},
	}
}
