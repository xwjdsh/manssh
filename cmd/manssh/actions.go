package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

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
	if ao.Config == nil {
		ao.Config = make(map[string]string)
	}

	if identityfile := c.String("identityfile"); identityfile != "" {
		ao.Config["identityfile"] = identityfile
	}

	if len(ao.Config) == 0 && ao.Connect == "" {
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
		NewAlias: c.String("rename"),
	}
	if kvConfig := c.Generic("config"); kvConfig != nil {
		uo.Config = kvConfig.(*kvFlag).m
	}
	if uo.Config == nil {
		uo.Config = make(map[string]string)
	}

	if identityfile := c.String("identityfile"); identityfile != "" || c.IsSet("identityfile") {
		uo.Config["identityfile"] = identityfile
	}
	if !uo.Valid() {
		return cli.NewExitError("the update option is invalid", 1)
	}

	host, err := manssh.Update(path, uo)
	if err != nil {
		fmt.Printf(utils.ErrorFlag)
		return cli.NewExitError(err, 1)
	}

	fmt.Printf("%s updated successfully\n\n", utils.SuccessFlag)
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
	fmt.Printf("%s deleted successfully\n\n", utils.SuccessFlag)
	printHosts(c.Bool("path"), hosts)
	return nil
}

func backupCmd(c *cli.Context) error {
	if err := utils.ArgumentsCheck(c.NArg(), 1, 1); err != nil {
		return printErrorWithHelp(c, err)
	}
	backupPath := c.Args().First()
	if err := os.MkdirAll(backupPath, os.ModePerm); err != nil {
		return cli.NewExitError(err, 1)
	}

	paths, err := manssh.GetFilePaths(path)
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	pathDir := filepath.Dir(path)
	for _, p := range paths {
		bp := backupPath
		if p != path && strings.HasPrefix(p, pathDir) {
			bp = filepath.Join(bp, strings.Replace(p, pathDir, "", 1))
			if err := os.MkdirAll(filepath.Dir(bp), os.ModePerm); err != nil {
				return cli.NewExitError(err, 1)
			}
		}
		if err := exec.Command("cp", p, bp).Run(); err != nil {
			return cli.NewExitError(err, 1)
		}
	}
	fmt.Printf("%s backup ssh config to [%s] successfully\n", utils.SuccessFlag, backupPath)
	return nil
}

func open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

func webCmd(c *cli.Context) error {
	addr := fmt.Sprintf("%s:%d", c.String("bind"), c.Int("port"))
	url := fmt.Sprintf("http://%s", addr)
	fmt.Printf("Running at: %s\n", url)
	go open(url)
	return manssh.WebServe(path, addr, c.Bool("cors"))
}
