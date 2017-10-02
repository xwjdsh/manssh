package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = name
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
		{Name: "update", Usage: "update existing ssh config record", Action: update},
		{Name: "delete", Usage: "delete existing ssh config record", Action: delete},
		{Name: "rename", Usage: "rename existing ssh config record", Action: rename},
		{Name: "export", Usage: "export ssh config record", Action: export},
	}
}
