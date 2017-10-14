package main

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/urfave/cli"
)

var (
	version = "master"
)

const (
	usage = "Manage your ssh alias configs easily"
)

var (
	whiteBoldColor   = color.New(color.FgWhite, color.Bold)
	yellowBoldColor  = color.New(color.FgYellow, color.Bold)
	magentaBoldColor = color.New(color.FgMagenta, color.Bold)
	successColor     = color.New(color.BgGreen, color.FgWhite)
	errorColor       = color.New(color.BgRed, color.FgWhite)
)

func formatConnect(user, hostname, port string) string {
	return fmt.Sprintf("%s@%s:%s", user, hostname, port)
}

// format is [user@]host[:port]
func parseConnct(connect string) (string, string, string) {
	var u, hostname, port string
	port = "22"
	currentUser, _ := user.Current()
	u = currentUser.Name

	hs := strings.Split(connect, "@")
	hostname = hs[0]
	if len(hs) > 1 {
		if hs[0] != "" {
			u = hs[0]
		}
		hostname = hs[1]
	}
	hss := strings.Split(hostname, ":")
	hostname = hss[0]
	if len(hss) > 1 {
		if _, err := strconv.Atoi(hss[1]); err == nil {
			port = hss[1]
		}
	}
	return u, hostname, port
}

func printSuccessFlag() {
	successColor.Printf("%-9s", " success")
}

func printErrorFlag() {
	errorColor.Printf("%-7s", " error")
}

func printErrorWithHelp(c *cli.Context, err error) error {
	cli.ShowSubcommandHelp(c)
	fmt.Println()
	printErrorFlag()
	return cli.NewExitError(err, 1)
}

func printHost(host *hostConfig) {
	yellowBoldColor.Printf("\t%s", host.aliases)
	fmt.Printf(" -> %s\n", host.connect)
	for k, v := range host.config {
		fmt.Printf("\t\t%s = %s\n", k, v)
	}
	fmt.Println()
}

func printGlobalConfig(config map[string]string) {
	magentaBoldColor.Printf("\t (*) Global Configs\n")
	for k, v := range config {
		fmt.Printf("\t\t%s = %s\n", k, v)
	}
	fmt.Println()
}

func argumentsCheck(c *cli.Context, min, max int) error {
	argCount := c.NArg()
	var err error
	if min > 0 && argCount < min {
		err = errors.New("too few arguments.")
	}
	if max > 0 && argCount > max {
		err = errors.New("too many arguments.")
	}
	return err
	if err != nil {
	}
	return nil
}

func query(values, keys []string) bool {
	for _, key := range keys {
		if !contains(values, key) {
			return false
		}
	}
	return true
}

func contains(values []string, key string) bool {
	for _, value := range values {
		if strings.Contains(value, key) {
			return true
		}
	}
	return false
}

func getHomeDir() string {
	user, err := user.Current()
	if nil == err && user.HomeDir != "" {
		return user.HomeDir
	}
	return os.Getenv("HOME")
}
