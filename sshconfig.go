package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/fatih/color"
	"github.com/mikkeloscar/sshconfig"
	"github.com/urfave/cli"
)

var (
	path            string
	whiteBoldColor  = color.New(color.FgWhite, color.Bold)
	yellowBoldColor = color.New(color.FgYellow, color.Bold)
)

func formatHost(host *sshconfig.SSHHost) string {
	return fmt.Sprintf("%s@%s:%d", host.User, host.HostName, host.Port)
}

func list(c *cli.Context) error {
	hosts, _ := sshconfig.ParseSSHConfig(path)
	if len(c.Args()) > 0 {
		searchHosts := []*sshconfig.SSHHost{}
		for _, host := range hosts {
			values := []string{host.HostName, host.User, fmt.Sprintf("%d", host.Port)}
			values = append(values, host.Host...)
			if query(values, c.Args()) {
				searchHosts = append(searchHosts, host)
			}
		}
		hosts = searchHosts
	}
	whiteBoldColor.Printf("Display %d records:\n\n", len(hosts))
	for _, host := range hosts {
		yellowBoldColor.Printf("    %s", strings.Join(host.Host, " "))
		fmt.Printf(" -> %s\n\n", formatHost(host))
	}
	return nil
}

func add(c *cli.Context) error {
	fmt.Println("add command")
	return nil
}

func update(c *cli.Context) error {
	fmt.Println("update command")
	return nil
}

func delete(c *cli.Context) error {
	fmt.Println("delete command")
	return nil
}

func rename(c *cli.Context) error {
	fmt.Println("rename command")
	return nil
}

func export(c *cli.Context) error {
	if len(c.Args()) == 0 {
		cli.ShowSubcommandHelp(c)
		fmt.Println()
		return cli.NewExitError("arguments is missing", 1)
	}

	if len(c.Args()) > 1 {
		cli.ShowSubcommandHelp(c)
		fmt.Println()
		return cli.NewExitError("too many arguments", 1)
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	err = ioutil.WriteFile(c.Args().First(), data, 0644)
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	return nil
}
