package main

import (
	"io/ioutil"

	"github.com/urfave/cli"
)

var (
	path string
)

func listAction(c *cli.Context) error {
	hosts, globalConfig := listHost(c.Args()...)
	printSuccessFlag()
	whiteBoldColor.Printf("display %d records.\n\n", len(hosts))
	for _, host := range hosts {
		printHost(host)
	}
	if len(globalConfig) > 0 {
		printGlobalConfig(globalConfig)
	}
	return nil
}

func addAction(c *cli.Context) error {
	if err := argumentsCheck(c.NArg(), 2, 2); err != nil {
		return printErrorWithHelp(c, err)
	}
	alias := c.Args().Get(0)
	connect := c.Args().Get(1)
	host := &hostConfig{
		aliases: alias,
		connect: connect,
	}
	if kvConfig := c.Generic("config"); kvConfig != nil {
		host.config = kvConfig.(*kvFlag).m
	}
	if err := addHost(host); err != nil {
		printErrorFlag()
		return cli.NewExitError(err, 1)
	}
	printSuccessFlag()
	whiteBoldColor.Printf("alias[%s] added successfully.\n\n", alias)
	printHost(host)
	return nil
}

func updateAction(c *cli.Context) error {
	if err := argumentsCheck(c.NArg(), 1, 2); err != nil {
		return printErrorWithHelp(c, err)
	}
	alias := c.Args().Get(0)
	connect := c.Args().Get(1)
	host := &hostConfig{
		aliases: alias,
		connect: connect,
	}
	if kvConfig := c.Generic("config"); kvConfig != nil {
		host.config = kvConfig.(*kvFlag).m
	}

	if err := updateHost(host, c.String("rename")); err != nil {
		printErrorFlag()
		return cli.NewExitError(err, 1)
	}

	printSuccessFlag()
	whiteBoldColor.Printf("alias[%s] updated successfully.\n\n", alias)
	printHost(host)
	return nil
}

func deleteAction(c *cli.Context) error {
	if err := argumentsCheck(c.NArg(), 1, -1); err != nil {
		return printErrorWithHelp(c, err)
	}
	if err := deleteHost(c.Args()...); err != nil {
		printErrorFlag()
		return cli.NewExitError(err, 1)
	}
	printSuccessFlag()
	whiteBoldColor.Printf("deleted %d records.\n", len(c.Args()))
	return nil
}

func backupAction(c *cli.Context) error {
	if err := argumentsCheck(c.NArg(), 1, 1); err != nil {
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
	whiteBoldColor.Printf("backup ssh config to [%s] successfully.", backupPath)
	return nil
}
