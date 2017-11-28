package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/urfave/cli"
	"github.com/xwjdsh/manssh"
)

var (
	messageStyle = color.New(color.FgWhite, color.Bold)
	aliasStyle   = color.New(color.FgYellow, color.Bold)

	globalStyle    = color.New(color.FgMagenta, color.Bold)
	globalKeyStyle = color.New(color.FgCyan, color.Bold)

	successStyle = color.New(color.FgGreen)
	errorStyle   = color.New(color.FgRed)
)

func printSuccessFlag() {
	successStyle.Printf("%s", "\u2714  ")
}

func printErrorFlag() {
	errorStyle.Printf("%s", "\u2716  ")
}

func printErrorWithHelp(c *cli.Context, err error) error {
	cli.ShowSubcommandHelp(c)
	fmt.Println()
	printErrorFlag()
	return cli.NewExitError(err, 1)
}

func printHosts(hosts []*manssh.HostConfig) {
	var global *manssh.HostConfig
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
	messageStyle.Printf(format, a...)
}

func printHost(host *manssh.HostConfig) {
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
