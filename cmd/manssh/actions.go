package main

import (
	"errors"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/urfave/cli"
	"github.com/xwjdsh/manssh"
)

var (
	path string
)

func list(c *cli.Context) error {
	hosts, err := manssh.List(path, c.Args(), c.Bool("ignorecase"))
	if err != nil {
		printErrorFlag()
		return cli.NewExitError(err, 1)
	}
	printSuccessFlag()
	printMessage("Listing %d records.\n\n", len(hosts))
	printHosts(hosts)
	return nil
}

func add(c *cli.Context) error {
	// Check arguments count
	if err := manssh.ArgumentsCheck(c.NArg(), 1, 2); err != nil {
		return printErrorWithHelp(c, err)
	}
	host := &manssh.HostConfig{
		Aliases: c.Args().Get(0),
		Connect: c.Args().Get(1),
	}
	if kvConfig := c.Generic("config"); kvConfig != nil {
		host.Config = kvConfig.(*kvFlag).m
	}
	if identityfile := c.String("identityfile"); identityfile != "" {
		if host.Config == nil {
			host.Config = map[string]string{}
		}
		host.Config[manssh.IdentityFile] = identityfile
	}

	if host.Config == nil && host.Connect == "" {
		return printErrorWithHelp(c, errors.New("param error"))
	}

	addPath := c.String("path")
	if addPath != "" {
		var err error
		addPath, err = filepath.Abs(addPath)
		if err != nil {
			printErrorFlag()
			return cli.NewExitError(err, 1)
		}
	}

	if err := manssh.Add(path, host, addPath); err != nil {
		printErrorFlag()
		return cli.NewExitError(err, 1)
	}
	printSuccessFlag()
	printMessage("alias[%s] added successfully.\n\n", host.Aliases)
	printHost(host)
	return nil
}

func update(c *cli.Context) error {
	if err := manssh.ArgumentsCheck(c.NArg(), 1, 2); err != nil {
		return printErrorWithHelp(c, err)
	}
	host := &manssh.HostConfig{
		Aliases: c.Args().Get(0),
		Connect: c.Args().Get(1),
	}
	if kvConfig := c.Generic("config"); kvConfig != nil {
		host.Config = kvConfig.(*kvFlag).m
	}
	c.FlagNames()
	if identityfile := c.String("identityfile"); identityfile != "" || c.IsSet("identityfile") {
		if host.Config == nil {
			host.Config = map[string]string{}
		}
		host.Config[manssh.IdentityFile] = identityfile
	}

	if err := manssh.Update(path, host, c.String("rename")); err != nil {
		printErrorFlag()
		return cli.NewExitError(err, 1)
	}

	printSuccessFlag()
	printMessage("alias[%s] updated successfully.\n\n", host.Aliases)
	printHost(host)
	return nil
}

func delete(c *cli.Context) error {
	if err := manssh.ArgumentsCheck(c.NArg(), 1, -1); err != nil {
		return printErrorWithHelp(c, err)
	}
	if err := manssh.Delete(path, c.Args()...); err != nil {
		printErrorFlag()
		return cli.NewExitError(err, 1)
	}
	printSuccessFlag()
	printMessage("alias[%s] deleted successfully.\n", strings.Join(c.Args(), ","))
	return nil
}

func backup(c *cli.Context) error {
	if err := manssh.ArgumentsCheck(c.NArg(), 1, 1); err != nil {
		return printErrorWithHelp(c, err)
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
	printMessage("backup ssh config to [%s] successfully.", backupPath)
	return nil
}
