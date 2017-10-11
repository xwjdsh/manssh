package main

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/howeyc/gopass"
	"github.com/mikkeloscar/sshconfig"
	"github.com/urfave/cli"
)

var (
	version = "master"
)

const (
	usage = "Manage your ssh alias configs easily"
)

var (
	whiteBoldColor  = color.New(color.FgWhite, color.Bold)
	yellowBoldColor = color.New(color.FgYellow, color.Bold)
	successColor    = color.New(color.BgGreen, color.FgWhite)
	errorColor      = color.New(color.BgRed, color.FgWhite)
)

// parsePath resolve path and get the real file path, if config file is a symbol link
func parsePath(path string) string {
	fileInfo, err := os.Lstat(path)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if fileInfo.Mode()&os.ModeSymlink != 0 {
		originFile, err := os.Readlink(path)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		return originFile
	}
	return path
}

func saveHosts(hosts []*sshconfig.SSHHost) error {
	var buffer bytes.Buffer
	for _, host := range hosts {
		buffer.WriteString(fmt.Sprintf("Host %s\n", strings.Join(host.Host, " ")))
		buffer.WriteString(fmt.Sprintf("    User %s\n", host.User))
		buffer.WriteString(fmt.Sprintf("    HostName %s\n", host.HostName))
		buffer.WriteString(fmt.Sprintf("    Port %d\n", host.Port))
		if host.IdentityFile != "" {
			buffer.WriteString(fmt.Sprintf("    IdentityFile %s\n", host.IdentityFile))
		}
		if host.ProxyCommand != "" {
			buffer.WriteString(fmt.Sprintf("    ProxyCommand %s\n", host.ProxyCommand))
		}
	}
	if err := ioutil.WriteFile(path, buffer.Bytes(), 0644); err != nil {
		printErrorFlag()
		return cli.NewExitError(err, 1)
	}
	return nil
}

func parseHost(alias, hostStr string, originHost *sshconfig.SSHHost) *sshconfig.SSHHost {
	var host *sshconfig.SSHHost
	if originHost != nil {
		host = originHost
	} else {
		host = &sshconfig.SSHHost{
			Host: []string{alias},
		}
	}
	host.Port = 22
	u, _ := user.Current()
	host.User = u.Name

	hs := strings.Split(hostStr, "@")
	connectUrl := hs[0]
	if len(hs) > 1 {
		if hs[0] != "" {
			host.User = hs[0]
		}
		connectUrl = hs[1]
	}
	hss := strings.Split(connectUrl, ":")
	host.HostName = hss[0]
	if len(hss) > 1 {
		if port, err := strconv.Atoi(hss[1]); err == nil {
			host.Port = port
		}
	}
	return host
}

func getHostsMap(hosts []*sshconfig.SSHHost) map[string]*sshconfig.SSHHost {
	hostMap := map[string]*sshconfig.SSHHost{}
	for _, host := range hosts {
		for _, alias := range host.Host {
			hostMap[alias] = host
		}
	}
	return hostMap
}

func checkAlias(hosts []*sshconfig.SSHHost, expectExist bool, alias ...string) (*sshconfig.SSHHost, error) {
	hostMap := getHostsMap(hosts)
	for _, a := range alias {
		_, ok := hostMap[a]
		if !ok && expectExist {
			return nil, fmt.Errorf("ssh alias('%s') not found.", a)
		} else if ok && !expectExist {
			return nil, fmt.Errorf("ssh alias('%s') already exists.", a)
		}
	}
	var host *sshconfig.SSHHost
	if len(alias) == 1 {
		host = hostMap[alias[0]]
	}
	return host, nil
}

func formatHost(host *sshconfig.SSHHost) string {
	return fmt.Sprintf("%s@%s:%d", host.User, host.HostName, host.Port)
}

func printSuccessFlag() {
	successColor.Printf("%-9s", " success")
}

func printErrorFlag() {
	errorColor.Printf("%-6s", " error")
}

func printErrorWithHelp(c *cli.Context, err error) error {
	cli.ShowSubcommandHelp(c)
	fmt.Println()
	printErrorFlag()
	return cli.NewExitError(err, 1)
}

func printHost(host *sshconfig.SSHHost) {
	yellowBoldColor.Printf("\t%s", strings.Join(host.Host, " "))
	fmt.Printf(" -> %s\n", formatHost(host))
	if host.IdentityFile != "" {
		fmt.Printf("\t\tIdentityFile = %s\n", host.IdentityFile)
	}
	if host.ProxyCommand != "" {
		fmt.Printf("\t\tProxyCommand = %s\n", host.ProxyCommand)
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

func readPrivateKey(path string) ([]byte, error) {
	privateKey, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load identity: %v", err)
	}

	block, rest := pem.Decode(privateKey)
	if len(rest) > 0 {
		return nil, fmt.Errorf("extra data when decoding private key")
	}
	if !x509.IsEncryptedPEMBlock(block) {
		return privateKey, nil
	}

	passphrase := []byte(os.Getenv("IDENTITY_PASSPHRASE"))
	if len(passphrase) == 0 {
		fmt.Print("Enter passphrase: ")
		passphrase, err = gopass.GetPasswd()
		if err != nil {
			return nil, fmt.Errorf("couldn't read passphrase: %v", err)
		}
	}
	der, err := x509.DecryptPEMBlock(block, passphrase)
	if err != nil {
		return nil, fmt.Errorf("decrypt failed: %v", err)
	}

	privateKey = pem.EncodeToMemory(&pem.Block{
		Type:  block.Type,
		Bytes: der,
	})

	return privateKey, nil
}
