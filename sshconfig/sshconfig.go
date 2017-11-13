package sshconfig

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/kevinburke/ssh_config"
	"github.com/xwjdsh/manssh/utils"
)

const (
	// User is ssh config key
	User = "user"
	// Hostname is ssh config key
	Hostname = "hostname"
	// Port is ssh config key
	Port = "port"
)

// ParseConfig parse configs from ssh config file, return config object and alias map
func ParseConfig(path string) (*ssh_config.Config, map[string]*ssh_config.Host) {
	f, _ := os.OpenFile(path, os.O_APPEND|os.O_CREATE, 0600)
	cfg, _ := ssh_config.Decode(f)
	aliasMap := map[string]*ssh_config.Host{}
	for _, host := range cfg.Hosts {
		for _, pattern := range host.Patterns {
			aliasMap[pattern.String()] = host
		}
	}
	return cfg, aliasMap
}

// List ssh alias by optional keyword
func List(path string, keywords ...string) []*utils.HostConfig {
	cfg, _ := ParseConfig(path)
	hosts := []*utils.HostConfig{}

	// Convert to utils.HostConfig
	for _, host := range cfg.Hosts {
		aliases := []string{}
		for _, pattern := range host.Patterns {
			aliases = append(aliases, pattern.String())
		}
		h := &utils.HostConfig{
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
		if len(keywords) > 0 && !utils.Query(values, keywords) {
			continue
		}
		if !isGlobal {
			h.Connect = utils.FormatConnect(connectMap[User], connectMap[Hostname], connectMap[Port])
		}
		hosts = append(hosts, h)
	}

	// Format
	if len(cfg.Hosts) > 0 {
		ioutil.WriteFile(path, []byte(cfg.String()), 0644)
	}
	return hosts
}

// Add ssh host config to ssh config file
func Add(path string, host *utils.HostConfig) error {
	cfg, aliasMap := ParseConfig(path)
	isGlobal := host.Aliases == "*"
	// Alias should not exist. except "*" because it always existing
	if !isGlobal {
		if err := utils.CheckAlias(aliasMap, false, host.Aliases); err != nil {
			return err
		}
	}
	if host.Config == nil {
		host.Config = map[string]string{}
	}
	nodes := []ssh_config.Node{}
	checkKeyRepeat := map[string]bool{}

	// Parse connect string
	if host.Connect != "" {
		user, hostname, port := utils.ParseConnct(host.Connect)
		host.Connect = utils.FormatConnect(user, hostname, port)

		userKV := &ssh_config.KV{Key: User, Value: user}
		nodes = append(nodes, userKV.SetLeadingSpace(4))

		hostnameKV := &ssh_config.KV{Key: Hostname, Value: hostname}
		nodes = append(nodes, hostnameKV.SetLeadingSpace(4))

		portKV := &ssh_config.KV{Key: Port, Value: port}
		nodes = append(nodes, portKV.SetLeadingSpace(4))

		checkKeyRepeat[User] = true
		checkKeyRepeat[Hostname] = true
		checkKeyRepeat[Port] = true
	}

	// Get nodes and delete repeat config
	deleteKeys := []string{}
	for k, v := range host.Config {
		lk := strings.ToLower(k)
		if !checkKeyRepeat[lk] {
			node := &ssh_config.KV{Key: lk, Value: v}
			nodes = append(nodes, node.SetLeadingSpace(4))
			checkKeyRepeat[lk] = true
		} else {
			deleteKeys = append(deleteKeys, k)
		}
	}
	for _, deleteKey := range deleteKeys {
		delete(host.Config, deleteKey)
	}

	pattern, err := ssh_config.NewPattern(host.Aliases)
	if err != nil {
		return nil
	}
	newHost := &ssh_config.Host{
		Patterns: []*ssh_config.Pattern{pattern},
		Nodes:    nodes,
	}
	cfg.Hosts = append(cfg.Hosts, newHost)
	return ioutil.WriteFile(path, []byte(cfg.String()), 0644)
}

// Update existing record
func Update(path string, h *utils.HostConfig, newAlias string) error {
	cfg, aliasMap := ParseConfig(path)
	if err := utils.CheckAlias(aliasMap, true, h.Aliases); err != nil {
		return err
	}

	updateHost := aliasMap[h.Aliases]
	if newAlias != "" {
		// alias rename
		for _, pattern := range updateHost.Patterns {
			if pattern.String() == h.Aliases {
				pattern.SetStr(newAlias)
				h.Aliases = newAlias
			}
		}
	}

	updateKV := map[string]string{}
	if h.Connect != "" {
		user, hostname, port := utils.ParseConnct(h.Connect)
		updateKV[User] = user
		updateKV[Hostname] = hostname
		updateKV[Port] = port
	}

	if h.Config != nil {
		for k, v := range h.Config {
			updateKV[strings.ToLower(k)] = v
		}
	}
	h.Config = map[string]string{}
	connectMap := map[string]string{}
	if h.Aliases != "*" {
		connectMap[User] = ""
		connectMap[Hostname] = ""
		connectMap[Port] = ""
	}
	// update node
	for _, node := range updateHost.Nodes {
		switch t := node.(type) {
		case *ssh_config.KV:
			t.Key = strings.ToLower(t.Key)
			if value, ok := updateKV[t.Key]; ok {
				t.SetLeadingSpace(4)
				t.Value = value
				delete(updateKV, t.Key)
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
		kv := &ssh_config.KV{Key: k, Value: v}
		updateHost.Nodes = append(updateHost.Nodes, kv.SetLeadingSpace(4))
		if _, ok := connectMap[k]; ok {
			connectMap[k] = v
		} else {
			h.Config[k] = v
		}
	}
	h.Connect = utils.FormatConnect(connectMap[User], connectMap[Hostname], connectMap[Port])
	return ioutil.WriteFile(path, []byte(cfg.String()), 0644)
}

// Delete existing alias record
func Delete(path string, aliases ...string) error {
	cfg, aliasMap := ParseConfig(path)
	if err := utils.CheckAlias(aliasMap, true, aliases...); err != nil {
		return err
	}
	deleteAliasMap := map[string]bool{}
	for _, alias := range aliases {
		deleteAliasMap[alias] = true
	}
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
	return ioutil.WriteFile(path, []byte(cfg.String()), 0644)
}
