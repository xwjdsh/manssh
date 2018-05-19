package main

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/urfave/cli"
	"github.com/xwjdsh/manssh"
)

var (
	messageStyle = color.New(color.FgWhite, color.Bold)
	aliasStyle   = color.New(color.FgMagenta, color.Bold)

	globalStyle    = color.New(color.FgMagenta, color.Bold)
	globalKeyStyle = color.New(color.FgCyan, color.Bold)

	successStyle = color.New(color.FgGreen)
	errorStyle   = color.New(color.FgRed)
)

func printSuccessFlag() {
	successStyle.Print("✔ ")
}

func printErrorFlag() {
	errorStyle.Print("✗ ")
}

func printErrorWithHelp(c *cli.Context, err error) error {
	cli.ShowSubcommandHelp(c)
	fmt.Println()
	printErrorFlag()
	return cli.NewExitError(err, 1)
}

func printHosts(hosts []*manssh.HostConfig, showPath bool) {
	var global *manssh.HostConfig
	for _, host := range hosts {
		if host.Aliases == "*" {
			global = host
			continue
		}
		printHost(host, showPath)
	}
	if global != nil && global.Config != nil && len(global.Config) > 0 {
		printHost(global, showPath)
	}
}

func printMessage(format string, a ...interface{}) {
	messageStyle.Printf(format, a...)
}

func printHost(host *manssh.HostConfig, showPath bool) {
	isGlobal := host.Aliases == "*"
	if isGlobal {
		globalStyle.Printf("\t(*) Global Configs\n")
	} else {
		aliasStyle.Printf("\t%s", host.Aliases)
		if showPath && host.Path != "" {
			if homeDir := manssh.GetHomeDir(); strings.HasPrefix(host.Path, homeDir) {
				host.Path = strings.Replace(host.Path, homeDir, "~", 1)
			}
			fmt.Printf("(%s)", host.Path)
		}
		fmt.Printf(" -> %s\n", color.GreenString(host.Connect))
	}
	for k, v := range host.Config {
		if v == "" {
			continue
		}
		if isGlobal {
			globalKeyStyle.Printf("\t\t%s ", k)
			fmt.Printf("= %s\n", v)
		} else {
			fmt.Printf("\t\t%s = %s\n", k, v)
		}
	}
	fmt.Println()
}
