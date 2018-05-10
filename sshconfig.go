package manssh

import (
	"io/ioutil"
	"os"
	"os/user"
	"strings"

	"github.com/xwjdsh/ssh_config"
)

const (
	// User is a ssh config key
	User = "user"
	// Hostname is a ssh config key
	Hostname = "hostname"
	// Port is a ssh config key
	Port = "port"
	// IdentityFile is a ssh config key
	IdentityFile = "identityfile"
)

type sshConfigHost struct {
	*ssh_config.Host
	path string
}

func parseConfig(path string) (*ssh_config.Config, error) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE, 0600)
	if err != nil {
		return nil, err
	}
	return ssh_config.Decode(f)
}

func writeConfig(path string, cfg *ssh_config.Config) error {
	return ioutil.WriteFile(path, []byte(cfg.String()), 0644)
}

// ParseConfig parse configs from ssh config file, return config object and alias map
func ParseConfig(path string) (map[string]*ssh_config.Config, map[string]*sshConfigHost, error) {
	cfg, err := parseConfig(path)
	if err != nil {
		return nil, nil, err
	}

	aliasMap := map[string]*sshConfigHost{}
	configMap := map[string]*ssh_config.Config{path: cfg}

	addAlias := func(fp string, hosts ...*ssh_config.Host) {
		for _, host := range hosts {
			for _, pattern := range host.Patterns {
				aliasMap[pattern.String()] = &sshConfigHost{
					Host: host,
					path: fp,
				}
			}
		}
	}

	for _, host := range cfg.Hosts {
		for _, node := range host.Nodes {
			switch t := node.(type) {
			case *ssh_config.Include:
				for fp, config := range t.GetFiles() {
					configMap[fp] = config
					addAlias(fp, config.Hosts...)
				}
			}
		}
		addAlias(path, host)
	}
	return configMap, aliasMap, nil
}

// List ssh alias, filter by optional keyword
func List(path string, keywords []string, ignoreCase ...bool) ([]*HostConfig, error) {
	configMap, aliasMap, err := ParseConfig(path)
	if err != nil {
		return nil, err
	}
	hosts := []*HostConfig{}

	// Convert to HostConfig
	hostMap := map[*ssh_config.Host]bool{}
	for _, host := range aliasMap {
		if hostMap[host.Host] {
			continue
		}
		hostMap[host.Host] = true

		aliases := []string{}
		for _, pattern := range host.Patterns {
			aliases = append(aliases, pattern.String())
		}
		h := &HostConfig{
			Aliases: strings.Join(aliases, " "),
			Config:  map[string]string{},
		}
		isGlobal := h.Aliases == "*"
		connectMap := map[string]string{}
		if !isGlobal {
			connectMap[User] = ""
			connectMap[Hostname] = ""
			connectMap[Port] = ""
		}
		values := []string{}
		values = append(values, aliases...)

		for _, node := range host.Nodes {
			switch t := node.(type) {
			case *ssh_config.KV:
				t.Key = strings.ToLower(t.Key)
				values = append(values, t.Value)
				if _, ok := connectMap[t.Key]; ok {
					connectMap[t.Key] = t.Value
				} else {
					h.Config[t.Key] = t.Value
				}
				t.SetLeadingSpace(4)
			}
		}
		if isGlobal && len(h.Config) == 0 {
			continue
		}
		if keywords != nil && len(keywords) > 0 && !Query(values, keywords, len(ignoreCase) > 0 && ignoreCase[0]) {
			continue
		}
		if !isGlobal {
			if connectMap[User] == "" {
				user, _ := user.Current()
				connectMap[User] = user.Username
			}
			if connectMap[Port] == "" {
				connectMap[Port] = "22"
			}
			h.Connect = FormatConnect(connectMap[User], connectMap[Hostname], connectMap[Port])
		}
		hosts = append(hosts, h)
	}

	// Format
	for fp, cfg := range configMap {
		if len(cfg.Hosts) > 0 {
			if err := writeConfig(fp, cfg); err != nil {
				return nil, err
			}
		}
	}
	return hosts, nil
}

// Add ssh host config to ssh config file
func Add(path string, host *HostConfig, addPath string) error {
	configMap, aliasMap, err := ParseConfig(path)
	if err != nil {
		return err
	}

	cfg, ok := configMap[addPath]
	if !ok {
		cfg, err = parseConfig(addPath)
		if err != nil {
			return err
		}
	}

	isGlobal := host.Aliases == "*"
	// Alias should not exist. except "*" because it always existing
	if !isGlobal || (isGlobal && len(aliasMap["*"].Nodes) > 0) {
		if err := CheckAlias(aliasMap, false, host.Aliases); err != nil {
			return err
		}
	}
	if host.Config == nil {
		host.Config = map[string]string{}
	}
	nodes := []ssh_config.Node{}

	// Parse connect string
	if host.Connect != "" {
		user, hostname, port := ParseConnct(host.Connect)
		host.Connect = FormatConnect(user, hostname, port)

		userKV := &ssh_config.KV{Key: User, Value: user}
		nodes = append(nodes, userKV.SetLeadingSpace(4))

		hostnameKV := &ssh_config.KV{Key: Hostname, Value: hostname}
		nodes = append(nodes, hostnameKV.SetLeadingSpace(4))

		portKV := &ssh_config.KV{Key: Port, Value: port}
		nodes = append(nodes, portKV.SetLeadingSpace(4))
		delete(host.Config, User)
		delete(host.Config, Port)
		delete(host.Config, Hostname)
	}

	// Get nodes and delete repeat config
	for k, v := range host.Config {
		if v == "" {
			continue
		}
		lk := strings.ToLower(k)
		node := &ssh_config.KV{Key: lk, Value: v}
		nodes = append(nodes, node.SetLeadingSpace(4))
	}

	pattern, err := ssh_config.NewPattern(host.Aliases)
	if err != nil {
		return nil
	}
	patterns := []*ssh_config.Pattern{pattern}
	newHost := &ssh_config.Host{
		Patterns: patterns,
		Nodes:    nodes,
	}
	if !isGlobal {
		cfg.Hosts = append(cfg.Hosts, newHost)
	} else {
		*aliasMap["*"].Host = *newHost
	}

	return writeConfig(addPath, cfg)
}

// Update existing record
func Update(path string, h *HostConfig, newAlias string) error {
	configMap, aliasMap, err := ParseConfig(path)
	if err != nil {
		return err
	}
	if err := CheckAlias(aliasMap, true, h.Aliases); err != nil {
		return err
	}

	updateHost := aliasMap[h.Aliases]
	if newAlias != "" {
		// new alias should not exist
		if err := CheckAlias(aliasMap, false, newAlias); err != nil {
			return err
		}

		// rename alias
		for _, pattern := range updateHost.Patterns {
			if pattern.String() == h.Aliases {
				pattern.SetStr(newAlias)
				h.Aliases = newAlias
			}
		}
	}

	updateKV := map[string]string{}
	if h.Config != nil {
		for k, v := range h.Config {
			updateKV[strings.ToLower(k)] = v
		}
	}

	if h.Connect != "" {
		user, hostname, port := ParseConnct(h.Connect)
		updateKV[User] = user
		updateKV[Hostname] = hostname
		updateKV[Port] = port
	}
	h.Config = map[string]string{}
	connectMap := map[string]string{}
	if h.Aliases != "*" {
		connectMap[User] = ""
		connectMap[Hostname] = ""
		connectMap[Port] = ""
	}

	// update node
	for i := 0; i >= 0 && i < len(updateHost.Nodes); i++ {
		switch t := updateHost.Nodes[i].(type) {
		case *ssh_config.KV:
			t.Key = strings.ToLower(t.Key)
			if value, ok := updateKV[t.Key]; ok {
				delete(updateKV, t.Key)
				if value == "" {
					// Remove node
					updateHost.Nodes = append(updateHost.Nodes[:i], updateHost.Nodes[i+1:]...)
					i--
					continue
				}
				t.SetLeadingSpace(4)
				t.Value = value
			}
			if _, ok := connectMap[t.Key]; ok {
				connectMap[t.Key] = t.Value
			} else {
				h.Config[t.Key] = t.Value
			}
		}
	}
	// append new node
	for k, v := range updateKV {
		if v == "" {
			continue
		}
		kv := &ssh_config.KV{Key: k, Value: v}
		updateHost.Nodes = append(updateHost.Nodes, kv.SetLeadingSpace(4))
		if _, ok := connectMap[k]; ok {
			connectMap[k] = v
		} else {
			h.Config[k] = v
		}
	}
	if connectMap[User] == "" {
		user, _ := user.Current()
		connectMap[User] = user.Username
	}
	if connectMap[Port] == "" {
		connectMap[Port] = "22"
	}
	h.Connect = FormatConnect(connectMap[User], connectMap[Hostname], connectMap[Port])
	return writeConfig(updateHost.path, configMap[updateHost.path])
}

// Delete existing alias record
func Delete(path string, aliases ...string) error {
	configMap, aliasMap, err := ParseConfig(path)
	if err != nil {
		return err
	}
	if err := CheckAlias(aliasMap, true, aliases...); err != nil {
		return err
	}
	deleteAliasMap := map[string]bool{}
	for _, alias := range aliases {
		deleteAliasMap[alias] = true
	}

	pathMap := map[string]bool{}
	for alias, _ := range deleteAliasMap {
		fp := aliasMap[alias].path
		if pathMap[fp] {
			continue
		}
		pathMap[fp] = true
		cfg := configMap[fp]
		newHosts := []*ssh_config.Host{}
		for _, host := range cfg.Hosts {
			newPattern := []*ssh_config.Pattern{}
			for _, pattern := range host.Patterns {
				if !deleteAliasMap[pattern.String()] {
					newPattern = append(newPattern, pattern)
				}
			}
			if len(newPattern) > 0 {
				host.Patterns = newPattern
				newHosts = append(newHosts, host)
			}
		}
		cfg.Hosts = newHosts
		if err := writeConfig(fp, cfg); err != nil {
			return err
		}
	}
	return nil
}
