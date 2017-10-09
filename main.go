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
		cli.StringFlag{Name: "config, c", Value: fmt.Sprintf("%s/.ssh/config", os.Getenv("HOME")), Destination: &path},
	}
}

func commands() []cli.Command {
	return []cli.Command{
		{Name: "add", Usage: "add a new ssh alias record", Action: add,
			Flags: []cli.Flag{
				cli.StringFlag{Name: "file, f"},
				cli.StringFlag{Name: "proxy, p"},
			},
		},
		{Name: "list", Usage: "list all ssh alias records", Action: list},
		{
			Name: "update", Usage: "update existing ssh alias record", Action: update,
			Flags: []cli.Flag{
				cli.StringFlag{Name: "user, u"},
				cli.StringFlag{Name: "host, H"},
				cli.StringFlag{Name: "port, p"},
				cli.StringFlag{Name: "alias, a"},
				cli.StringFlag{Name: "file, f"},
				cli.StringFlag{Name: "proxy, P"},
			},
		},
		{Name: "delete", Usage: "delete existing ssh alias record", Action: delete},
		{Name: "backup", Usage: "backup ssh alias config records", Action: backup},
	}
}
