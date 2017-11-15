package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
	"github.com/xwjdsh/manssh"
)

var (
	version = "master"
)

func main() {
	app := cli.NewApp()
	app.Usage = "Manage your ssh alias configs easily"
	app.Version = version
	app.Flags = flags()
	app.Commands = commands()
	app.Run(os.Args)
}

func flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{Name: "file, f", Value: fmt.Sprintf("%s/.ssh/config", manssh.GetHomeDir()), Destination: &path},
	}
}
