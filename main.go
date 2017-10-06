package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Usage = useage
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
		{Name: "add", Usage: "add a record to ssh config", Action: add},
		{Name: "list", Usage: "list all existing ssh config record", Action: list},
		{
			Name: "update", Usage: "update existing ssh config record", Action: update,
			Flags: []cli.Flag{
				cli.StringFlag{Name: "user, u"},
				cli.StringFlag{Name: "host, H"},
				cli.StringFlag{Name: "port, p"},
				cli.StringFlag{Name: "alias, a"},
			},
		},
		{Name: "delete", Usage: "delete existing ssh config record", Action: delete},
		{Name: "backup", Usage: "backup ssh config record", Action: backup},
	}
}
