package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/howeyc/gopass"
	"github.com/mikkeloscar/sshconfig"
	"github.com/urfave/cli"
	"golang.org/x/crypto/ssh"
)

var (
	path string
	CMD  = "tell application \"Terminal\" to do script \"%s\" in selected tab of the front window"
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
		return printErrorWithHelp(c, err)
	}
	newAlias := c.Args().Get(0)
	hostStr := c.Args().Get(1)
	hosts, _ := sshconfig.ParseSSHConfig(path)
	if _, err := checkAlias(hosts, false, newAlias); err != nil {
		printErrorFlag()
		return cli.NewExitError(err, 1)
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
		return printErrorWithHelp(c, err)
	}
	alias := c.Args().Get(0)
	hostStr := c.Args().Get(1)
	hosts, _ := sshconfig.ParseSSHConfig(path)
	host, err := checkAlias(hosts, true, alias)
	if err != nil {
		printErrorFlag()
		return cli.NewExitError(err, 1)
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
		return printErrorWithHelp(c, err)
	}
	hosts, _ := sshconfig.ParseSSHConfig(path)
	if _, err := checkAlias(hosts, true, c.Args()...); err != nil {
		printErrorFlag()
		return cli.NewExitError(err, 1)
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
	whiteBoldColor.Printf("backup ssh config to ('%s') successfully.", backupPath)
	return nil
}

func open(c *cli.Context) error {
	if err := argumentsCheck(c, 1, 1); err != nil {
		return printErrorWithHelp(c, err)
	}

	alias := c.Args().Get(0)
	hosts, _ := sshconfig.ParseSSHConfig(path)
	if _, err := checkAlias(hosts, true, alias); err != nil {
		printErrorFlag()
		return cli.NewExitError(err, 1)
	}
	ssh := fmt.Sprintf("ssh %s", alias)
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

func run(c *cli.Context) error {
	if err := argumentsCheck(c, 2, 2); err != nil {
		return printErrorWithHelp(c, err)
	}
	alias := c.Args().Get(0)
	command := c.Args().Get(1)
	hosts, _ := sshconfig.ParseSSHConfig(path)
	host, err := checkAlias(hosts, true, alias)
	if err != nil {
		printErrorFlag()
		return cli.NewExitError(err, 1)
	}
	user := host.User
	if u := c.String("user"); u != "" {
		user = u
	}

	auth := []ssh.AuthMethod{}
	if c.Bool("password") {
		fmt.Print("Enter password: ")
		password, err := gopass.GetPasswd()
		if err != nil {
			printErrorFlag()
			return cli.NewExitError(err, 1)
		}
		auth = append(auth, ssh.Password(string(password)))
	} else {
		keyPath := fmt.Sprintf("%s/.ssh/id_rsa", os.Getenv("HOME"))
		if host.IdentityFile != "" {
			keyPath = host.IdentityFile
		}
		keyBytes, err := readPrivateKey(keyPath)
		if err != nil {
			printErrorFlag()
			return cli.NewExitError(err, 1)
		}
		key, err := ssh.ParsePrivateKey(keyBytes)
		if err != nil {
			printErrorFlag()
			return cli.NewExitError(err, 1)
		}
		auth = append(auth, ssh.PublicKeys(key))
	}
	clientConfig := &ssh.ClientConfig{
		User:            user,
		Auth:            auth,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}
	addr := fmt.Sprintf("%s:%d", host.HostName, host.Port)
	client, err := ssh.Dial("tcp", addr, clientConfig)
	if err != nil {
		printErrorFlag()
		return cli.NewExitError(err, 1)
	}

	// create session
	session, err := client.NewSession()
	if err != nil {
		printErrorFlag()
		return cli.NewExitError(err, 1)
	}
	defer session.Close()

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	if err := session.Run(command); err != nil {
		printErrorFlag()
		return cli.NewExitError(err, 1)
	}
	return nil
}
