package main

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/xwjdsh/manssh"
	"github.com/xwjdsh/manssh/utils"

	"github.com/urfave/cli"
)

var (
	path string
)

func listCmd(c *cli.Context) error {
	hosts, err := manssh.List(path, manssh.ListOption{
		Keywords:   c.Args(),
		IgnoreCase: c.Bool("ignorecase"),
	})
	if err != nil {
		fmt.Printf(utils.ErrorFlag)
		return cli.NewExitError(err, 1)
	}
	fmt.Printf("%s total records: %d\n\n", utils.SuccessFlag, len(hosts))
	printHosts(c.Bool("path"), hosts)
	return nil
}

func addCmd(c *cli.Context) error {
	// Check arguments count
	if err := utils.ArgumentsCheck(c.NArg(), 1, 2); err != nil {
		return printErrorWithHelp(c, err)
	}
	ao := &manssh.AddOption{
		Alias:   c.Args().Get(0),
		Connect: c.Args().Get(1),
		Path:    c.String("addpath"),
		Config:  map[string]string{},
	}
	if ao.Path != "" {
		var err error
		if ao.Path, err = filepath.Abs(ao.Path); err != nil {
			fmt.Printf(utils.ErrorFlag)
			return cli.NewExitError(err, 1)
		}
	}
	if kvConfig := c.Generic("config"); kvConfig != nil {
		ao.Config = kvConfig.(*kvFlag).m
	}

	if identityfile := c.String("identityfile"); identityfile != "" {
		ao.Config["identityfile"] = identityfile
	}

	if ao.Config == nil && ao.Connect == "" {
		return printErrorWithHelp(c, errors.New("param error"))
	}

	host, err := manssh.Add(path, ao)
	if err != nil {
		fmt.Printf(utils.ErrorFlag)
		return cli.NewExitError(err, 1)
	}
	fmt.Printf("%s added successfully\n", utils.SuccessFlag)
	if host != nil {
		fmt.Println()
		printHost(c.Bool("path"), host)
	}
	return nil
}

func updateCmd(c *cli.Context) error {
	if err := utils.ArgumentsCheck(c.NArg(), 1, 2); err != nil {
		return printErrorWithHelp(c, err)
	}
	uo := &manssh.UpdateOption{
		Alias:    c.Args().Get(0),
		Connect:  c.Args().Get(1),
		Config:   map[string]string{},
		NewAlias: c.String("rename"),
	}
	if kvConfig := c.Generic("config"); kvConfig != nil {
		uo.Config = kvConfig.(*kvFlag).m
	}

	if identityfile := c.String("identityfile"); identityfile != "" || c.IsSet("identityfile") {
		uo.Config["identityfile"] = identityfile
	}

	host, err := manssh.Update(path, uo)
	if err != nil {
		fmt.Printf(utils.ErrorFlag)
		return cli.NewExitError(err, 1)
	}

	fmt.Printf("%s updated successfully.\n\n", utils.SuccessFlag)
	printHost(c.Bool("path"), host)
	return nil
}

func deleteCmd(c *cli.Context) error {
	if err := utils.ArgumentsCheck(c.NArg(), 1, -1); err != nil {
		return printErrorWithHelp(c, err)
	}
	hosts, err := manssh.Delete(path, c.Args()...)
	if err != nil {
		fmt.Printf(utils.ErrorFlag)
		return cli.NewExitError(err, 1)
	}
	fmt.Printf("%s deleted successfully.\n\n", utils.SuccessFlag)
	printHosts(c.Bool("path"), hosts)
	return nil
}
