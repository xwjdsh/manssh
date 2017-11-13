package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/urfave/cli"
	"github.com/xwjdsh/manssh/utils"
)

var (
	whiteBoldColor   = color.New(color.FgWhite, color.Bold)
	yellowBoldColor  = color.New(color.FgYellow, color.Bold)
	magentaBoldColor = color.New(color.FgMagenta, color.Bold)

	successColor = color.New(color.BgGreen, color.FgWhite)
	errorColor   = color.New(color.BgRed, color.FgWhite)
)

func printSuccessFlag() {
	successColor.Printf("%-9s", " success")
}

func printErrorFlag() {
	errorColor.Printf("%-7s", " error")
}

func printErrorWithHelp(c *cli.Context, err error) error {
	cli.ShowSubcommandHelp(c)
	fmt.Println()
	printErrorFlag()
	return cli.NewExitError(err, 1)
}

func printHosts(hosts []*utils.HostConfig) {
	var global *utils.HostConfig
	for _, host := range hosts {
		if host.Aliases == "*" {
			global = host
		}
		printHost(host)
	}
	if global != nil && global.Config != nil && len(global.Config) > 0 {
		printHost(global)
	}
}

func printHost(host *utils.HostConfig) {
	if host.Aliases == "*" {
		magentaBoldColor.Printf("\t (*) Global Configs\n")
	} else {
		yellowBoldColor.Printf("\t%s", host.Aliases)
		fmt.Printf(" -> %s\n", host.Connect)
	}
	for k, v := range host.Config {
		fmt.Printf("\t\t%s = %s\n", k, v)
	}
	fmt.Println()
}
