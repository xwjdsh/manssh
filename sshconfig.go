package main

import (
	"bytes"
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
	successColor    = color.New(color.BgGreen, color.FgWhite)
	errorColor      = color.New(color.BgRed, color.FgWhite)
)

func saveHosts(hosts []*sshconfig.SSHHost) error {
	var buffer bytes.Buffer
	for _, host := range hosts {
		buffer.WriteString(fmt.Sprintf("Host %s\n", strings.Join(host.Host, " ")))
		buffer.WriteString(fmt.Sprintf("    user %s\n", host.User))
		buffer.WriteString(fmt.Sprintf("    hostname %s\n", host.HostName))
		buffer.WriteString(fmt.Sprintf("    port %d\n", host.Port))
	}
	return ioutil.WriteFile(path, buffer.Bytes(), 0644)
}

func formatHost(host *sshconfig.SSHHost) string {
	return fmt.Sprintf("%s@%s:%d", host.User, host.HostName, host.Port)
}

func printSuccessFlag() {
	successColor.Printf("%-10s", " success")
}

func printErrorFlag() {
	errorColor.Printf("%-8s", " error")
}

func printHost(host *sshconfig.SSHHost) {
	yellowBoldColor.Printf("    %s", strings.Join(host.Host, " "))
	fmt.Printf(" -> %s\n\n", formatHost(host))
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
	printSuccessFlag()
	whiteBoldColor.Printf("Display %d records:\n\n", len(hosts))
	for _, host := range hosts {
		printHost(host)
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
	if err := argumentsCheck(c.Args(), 2, 2); err != nil {
		cli.ShowSubcommandHelp(c)
		fmt.Println()
		printErrorFlag()
		return cli.NewExitError(err, 1)
	}

	hosts, _ := sshconfig.ParseSSHConfig(path)
	hostMap := map[string]*sshconfig.SSHHost{}
	for _, host := range hosts {
		for _, alias := range host.Host {
			hostMap[alias] = host
		}
	}
	oldName := c.Args().Get(0)
	newName := c.Args().Get(1)
	if _, ok := hostMap[oldName]; !ok {
		printErrorFlag()
		return cli.NewExitError("old ssh alias not found", 1)
	}
	if _, ok := hostMap[newName]; ok {
		printErrorFlag()
		return cli.NewExitError("new ssh alias already exists", 1)
	}
	host := hostMap[oldName]
	for i, name := range host.Host {
		if name == oldName {
			host.Host[i] = newName
			break
		}
	}
	if err := saveHosts(hosts); err != nil {
		printErrorFlag()
		return cli.NewExitError(err, 1)
	}
	printSuccessFlag()
	whiteBoldColor.Printf("Rename from '%s' to '%s'\n\n", oldName, newName)
	printHost(host)
	return nil
}

func backup(c *cli.Context) error {
	if err := argumentsCheck(c.Args(), 1, 1); err != nil {
		cli.ShowSubcommandHelp(c)
		fmt.Println()
		printErrorFlag()
		return cli.NewExitError(err, 1)
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		printErrorFlag()
		return cli.NewExitError(err, 1)
	}
	backupPath := c.Args().First()
	err = ioutil.WriteFile(backupPath, data, 0644)
	if err != nil {
		printErrorFlag()
		return cli.NewExitError(err, 1)
	}
	printSuccessFlag()
	fmt.Printf("backup ssh config to '%s' success", backupPath)
	return nil
}
