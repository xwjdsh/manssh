package main

import (
	"errors"
	"io/ioutil"
	"strings"

	"github.com/urfave/cli"
	"github.com/xwjdsh/manssh/sshconfig"
	"github.com/xwjdsh/manssh/utils"
)

var (
	path string
)

func list(c *cli.Context) error {
	hosts := sshconfig.List(path, c.Args()...)
	printSuccessFlag()
	printMessage("Listing %d records.\n\n", len(hosts))
	printHosts(hosts)
	return nil
}

func add(c *cli.Context) error {
	// Check arguments count
	if err := utils.ArgumentsCheck(c.NArg(), 1, 2); err != nil {
		return printErrorWithHelp(c, err)
	}
	host := &utils.HostConfig{
		Aliases: c.Args().Get(0),
		Connect: c.Args().Get(1),
	}
	if kvConfig := c.Generic("config"); kvConfig != nil {
		host.Config = kvConfig.(*kvFlag).m
	}

	if host.Config == nil && host.Connect == "" {
		return printErrorWithHelp(c, errors.New("param error"))
	}

	if err := sshconfig.Add(path, host); err != nil {
		printErrorFlag()
		return cli.NewExitError(err, 1)
	}
	printSuccessFlag()
	printMessage("alias[%s] added successfully.\n\n", host.Aliases)
	printHost(host)
	return nil
}

func update(c *cli.Context) error {
	if err := utils.ArgumentsCheck(c.NArg(), 1, 2); err != nil {
		return printErrorWithHelp(c, err)
	}
	host := &utils.HostConfig{
		Aliases: c.Args().Get(0),
		Connect: c.Args().Get(1),
	}
	if kvConfig := c.Generic("config"); kvConfig != nil {
		host.Config = kvConfig.(*kvFlag).m
	}

	if err := sshconfig.Update(path, host, c.String("rename")); err != nil {
		printErrorFlag()
		return cli.NewExitError(err, 1)
	}

	printSuccessFlag()
	printMessage("alias[%s] updated successfully.\n\n", host.Aliases)
	printHost(host)
	return nil
}

func delete(c *cli.Context) error {
	if err := utils.ArgumentsCheck(c.NArg(), 1, -1); err != nil {
		return printErrorWithHelp(c, err)
	}
	if err := sshconfig.Delete(path, c.Args()...); err != nil {
		printErrorFlag()
		return cli.NewExitError(err, 1)
	}
	printSuccessFlag()
	printMessage("alias[%s] deleted successfully.\n", strings.Join(c.Args(), ","))
	return nil
}

func backup(c *cli.Context) error {
	if err := utils.ArgumentsCheck(c.NArg(), 1, 1); err != nil {
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
