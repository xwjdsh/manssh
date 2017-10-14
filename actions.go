package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/urfave/cli"
)

var (
	path string
	CMD  = "tell application \"Terminal\" to do script \"%s\" in selected tab of the front window"
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
	if err := argumentsCheck(c, 2, 2); err != nil {
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
	whiteBoldColor.Printf("alias(%s) added successfully.\n\n", alias)
	printHost(host)
	return nil
}

func updateAction(c *cli.Context) error {
	if err := argumentsCheck(c, 1, 2); err != nil {
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
	whiteBoldColor.Printf("alias(%s) updated successfully.\n\n", alias)
	printHost(host)
	return nil
}

func deleteAction(c *cli.Context) error {
	if err := argumentsCheck(c, 1, -1); err != nil {
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
	if err := argumentsCheck(c, 1, 1); err != nil {
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
	whiteBoldColor.Printf("backup ssh config to (%s) successfully.", backupPath)
	return nil
}

func openAction(c *cli.Context) error {
	if err := argumentsCheck(c, 1, 1); err != nil {
		return printErrorWithHelp(c, err)
	}
	addr := c.Args().Get(0)
	ssh := fmt.Sprintf("ssh %s", addr)
	sshCMD := exec.Command("/usr/bin/osascript", "-e", "tell application \"Terminal\" to activate", "-e", "tell application \"System Events\" to tell process \"Terminal\" to keystroke \"t\" using command down", "-e", fmt.Sprintf(CMD, ssh))
	sshCMD.Env = os.Environ()
	output, err := sshCMD.CombinedOutput()

	if err != nil {
		printErrorFlag()
		return cli.NewExitError(err, 1)
	}

	printSuccessFlag()
	whiteBoldColor.Println("out: ", string(output))
	return nil
}

func runAction(c *cli.Context) error {
	if err := argumentsCheck(c, 2, 2); err != nil {
		return printErrorWithHelp(c, err)
	}
	alias := c.Args().Get(0)
	command := c.Args().Get(1)

	userAlias := strings.Split(alias, "@")
	tmpUser, tmpAlias := "", alias
	if len(userAlias) == 2 {
		tmpUser, tmpAlias = userAlias[0], userAlias[1]
	}

	exists, user, hostname, port, identityfile := getHostConnect(tmpAlias)
	if !exists {
		user, hostname, port = parseConnct(alias)
	} else if tmpUser != "" {
		user = tmpUser
	}
	log.Println(user, hostname, port)
	session, err := createSession(c.Bool("password"), user, hostname, port, identityfile)
	if err != nil {
		printErrorFlag()
		return cli.NewExitError(err, 1)
	}
	if err := executeCommand(session, command); err != nil {
		printErrorFlag()
		return cli.NewExitError(err, 1)
	}

	return nil
}
