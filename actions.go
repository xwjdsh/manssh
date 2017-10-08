package main

import (
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/mikkeloscar/sshconfig"
	"github.com/urfave/cli"
)

var (
	path string
)

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
	whiteBoldColor.Printf("display '%d' records.\n\n", len(hosts))
	for _, host := range hosts {
		printHost(host)
	}
	return nil
}

func add(c *cli.Context) error {
	if err := argumentsCheck(c, 2, 2); err != nil {
		return err
	}
	newAlias := c.Args().Get(0)
	hostStr := c.Args().Get(1)
	hosts, _ := sshconfig.ParseSSHConfig(path)
	hostMap := getHostsMap(hosts)
	if _, ok := hostMap[newAlias]; ok {
		printErrorFlag()
		return cli.NewExitError(fmt.Sprintf("ssh alias('%s') already exists.", newAlias), 1)
	}
	host := parseHost(newAlias, hostStr, nil)
	if identityFile := c.String("file"); identityFile != "" {
		host.IdentityFile = identityFile
	}
	if proxyCommand := c.String("proxy"); proxyCommand != "" {
		host.ProxyCommand = proxyCommand
	}
	hosts = append(hosts, host)
	if err := saveHosts(hosts); err != nil {
		return err
	}
	printSuccessFlag()
	whiteBoldColor.Printf("ssh alias('%s') added successfully.\n\n", newAlias)
	printHost(host)
	return nil
}

func update(c *cli.Context) error {
	if err := argumentsCheck(c, 1, 2); err != nil {
		return err
	}
	alias := c.Args().Get(0)
	hostStr := c.Args().Get(1)
	hosts, _ := sshconfig.ParseSSHConfig(path)
	hostMap := getHostsMap(hosts)
	host, ok := hostMap[alias]
	if !ok {
		printErrorFlag()
		return cli.NewExitError(fmt.Sprintf("ssh alias('%s') not found.", alias), 1)
	}
	newUser, newHostname, newPort, newAlias := c.String("user"), c.String("host"), c.String("port"), c.String("alias")
	newIdentityFile, newProxy := c.String("file"), c.String("proxy")
	if c.NArg() == 1 && newUser == "" && newHostname == "" && newPort == "" && newAlias == "" && newIdentityFile == "" && newProxy == "" {
		printErrorFlag()
		return cli.NewExitError("too few arguments.", 1)
	}
	if hostStr != "" {
		parseHost(alias, hostStr, host)
	}
	if newUser != "" {
		host.User = newUser
	}
	if newHostname != "" {
		host.HostName = newHostname
	}
	if newPort != "" {
		if p, err := strconv.Atoi(newPort); err != nil {
			host.Port = p
		}
	}
	if newAlias != "" {
		for i, name := range host.Host {
			if name == alias {
				host.Host[i] = newAlias
				break
			}
		}
	}
	if newIdentityFile != "" {
		host.IdentityFile = newIdentityFile
	}
	if newProxy != "" {
		host.ProxyCommand = newProxy
	}
	if err := saveHosts(hosts); err != nil {
		return err
	}
	printSuccessFlag()
	whiteBoldColor.Printf("ssh alias('%s') updated successfully.\n\n", alias)
	printHost(host)
	return nil
}

func delete(c *cli.Context) error {
	if err := argumentsCheck(c, 1, -1); err != nil {
		return err
	}
	hosts, _ := sshconfig.ParseSSHConfig(path)
	hostMap := getHostsMap(hosts)
	for _, alias := range c.Args() {
		if _, ok := hostMap[alias]; !ok {
			printErrorFlag()
			return cli.NewExitError(fmt.Sprintf("ssh alias('%s') not found.", alias), 1)
		}
	}
	newHosts := []*sshconfig.SSHHost{}
	for _, host := range hosts {
		newAlias := []string{}
		for _, hostAlias := range host.Host {
			isDelete := false
			for _, deleteAlias := range c.Args() {
				if hostAlias == deleteAlias {
					isDelete = true
					break
				}
			}
			if !isDelete {
				newAlias = append(newAlias, hostAlias)
			}
		}
		host.Host = newAlias
		if len(host.Host) > 0 {
			newHosts = append(newHosts, host)
		}
	}
	if err := saveHosts(newHosts); err != nil {
		return err
	}
	printSuccessFlag()
	whiteBoldColor.Printf("deleted '%d' records.\n", len(c.Args()))
	return nil
}

func backup(c *cli.Context) error {
	if err := argumentsCheck(c, 1, 1); err != nil {
		return err
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
	whiteBoldColor.Printf("backup ssh config to ('%s') successfully.", backupPath)
	return nil
}
