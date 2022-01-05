package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"

	"github.com/xwjdsh/manssh/utils"
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
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{Name: "file, f", Value: fmt.Sprintf("%s/.ssh/config", utils.GetHomeDir()), Destination: &path},
	}
}
