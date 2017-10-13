package main

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net"
	"os"

	"github.com/howeyc/gopass"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"golang.org/x/crypto/ssh/terminal"
)

var keyPath = fmt.Sprintf("%s/.ssh/id_rsa", os.Getenv("HOME"))

func getAuthMethod(passwordAuth bool) (ssh.AuthMethod, error) {
	var auth ssh.AuthMethod
	var err error
	if passwordAuth {
		fmt.Print("Enter password: ")
		password, err := gopass.GetPasswd()
		if err != nil {
			return nil, err
		}
		auth = ssh.Password(string(password))
	} else {
		auth, err = publicKeyFile(keyPath)
	}
	return auth, err
}

func sshAgent() ssh.AuthMethod {
	if sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		return ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers)
	}
	return nil
}

func publicKeyFile(file string) (ssh.AuthMethod, error) {
	buffer, err := readPrivateKey(file)
	if err != nil {
		return nil, err
	}
	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil, err
	}
	return ssh.PublicKeys(key), nil
}

func readPrivateKey(file string) ([]byte, error) {
	privateKey, err := ioutil.ReadFile(file)
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

func createSession(passwordAuth bool, user, hostname, port string) (*ssh.Session, error) {
	auth, err := getAuthMethod(passwordAuth)
	if err != nil {
		return nil, err
	}
	var sshConfig = &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{auth},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	connection, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", hostname, port), sshConfig)
	if err != nil {
		return nil, err
	}
	return connection.NewSession()
}

func executeCommand(session *ssh.Session, command string) error {
	defer session.Close()

	fd := int(os.Stdin.Fd())
	oldState, err := terminal.MakeRaw(fd)
	if err != nil {
		return err
	}
	defer terminal.Restore(fd, oldState)

	// excute command
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	termWidth, termHeight, err := terminal.GetSize(fd)
	if err != nil {
		return err
	}

	// Set up terminal modes
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // enable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	// Request pseudo terminal
	if err := session.RequestPty("xterm-256color", termHeight, termWidth, modes); err != nil {
		return err
	}
	return session.Run(command)
}
