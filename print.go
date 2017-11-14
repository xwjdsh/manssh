package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/urfave/cli"
	"github.com/xwjdsh/manssh/utils"
)

var (
	messageStyle   = color.New(color.FgWhite, color.Bold)
	aliasStyle     = color.New(color.FgYellow, color.Bold)
	globalStyle    = color.New(color.FgMagenta, color.Bold)
	globalKeyStyle = color.New(color.FgCyan, color.Bold)

	successStyle = color.New(color.BgGreen, color.FgWhite)
	errorStyle   = color.New(color.BgRed, color.FgWhite)
)

func printSuccessFlag() {
	successStyle.Printf("%-9s", " success")
}

func printErrorFlag() {
	errorStyle.Printf("%-7s", " error")
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
			continue
		}
		printHost(host)
	}
	if global != nil && global.Config != nil && len(global.Config) > 0 {
		printHost(global)
	}
}

func printMessage(format string, a ...interface{}) {
	messageStyle.Printf(format, a)
}

func printHost(host *utils.HostConfig) {
	isGlobal := host.Aliases == "*"
	if isGlobal {
		globalStyle.Printf("\t(*) Global Configs\n")
	} else {
		aliasStyle.Printf("\t%s", host.Aliases)
		fmt.Printf(" -> %s\n", host.Connect)
	}
	for k, v := range host.Config {
		if isGlobal {
			globalKeyStyle.Printf("\t\t%s ", k)
			fmt.Printf("= %s\n", v)
		} else {
			fmt.Printf("\t\t%s = %s\n", k, v)
		}
	}
	fmt.Println()
}
