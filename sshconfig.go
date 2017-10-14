package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/kevinburke/ssh_config"
)

type hostConfig struct {
	aliases string
	connect string
	config  map[string]string
}

const (
	USER          = "user"
	HOSTNAME      = "hostname"
	PORT          = "port"
	IDENTITY_FILE = "identityfile"
)

func checkAlias(aliasMap map[string]*ssh_config.Host, expectExist bool, aliases ...string) error {
	for _, alias := range aliases {
		ok := aliasMap[alias] != nil
		if !ok && expectExist {
			return fmt.Errorf("alias(%s) not found.", alias)
		} else if ok && !expectExist {
			return fmt.Errorf("alias(%s) already exists.", alias)
		}
	}
	return nil
}

func parseConfig() (*ssh_config.Config, map[string]*ssh_config.Host) {
	f, _ := os.OpenFile(path, os.O_APPEND|os.O_CREATE, 0600)
	cfg, _ := ssh_config.Decode(f)
	aliasMap := map[string]*ssh_config.Host{}
	for _, host := range cfg.Hosts {
		for _, pattern := range host.Patterns {
			// exclude global config alias
			if alias := pattern.String(); alias != "*" {
				aliasMap[alias] = host
			}
		}
	}
	return cfg, aliasMap
}

func getHostConnect(alias string) (bool, string, string, string, string) {
	_, aliasMap := parseConfig()
	host := aliasMap[alias]
	if host == nil {
		return false, "", "", "", ""
	}
	connectMap := map[string]string{USER: "", HOSTNAME: "", PORT: "", IDENTITY_FILE: ""}
	for _, node := range host.Nodes {
		switch t := node.(type) {
		case *ssh_config.KV:
			lk := strings.ToLower(t.Key)
			if _, ok := connectMap[lk]; ok {
				connectMap[lk] = t.Value
			}
		}
	}
	return true, connectMap[USER], connectMap[HOSTNAME], connectMap[PORT], connectMap[IDENTITY_FILE]
}

func listHost(keywords ...string) ([]*hostConfig, map[string]string) {
	cfg, _ := parseConfig()
	hosts := []*hostConfig{}
	globalConfig := map[string]string{}
	for _, host := range cfg.Hosts {
		aliases := []string{}
		for _, pattern := range host.Patterns {
			aliases = append(aliases, pattern.String())
		}
		aliasesAgg := strings.Join(aliases, " ")
		configMap := map[string]string{}
		connectMap := map[string]string{USER: "", HOSTNAME: "", PORT: ""}
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
					configMap[t.Key] = t.Value
				}
				t.SetLeadingSpace(4)
			case *ssh_config.Include:
				// TODO handle include node
			}
		}
		if len(keywords) > 0 && !query(values, keywords) {
			continue
		}
		if aliasesAgg != "*" {
			host := &hostConfig{
				aliases: aliasesAgg,
				connect: formatConnect(connectMap["user"], connectMap["hostname"], connectMap["port"]),
				config:  configMap,
			}
			hosts = append(hosts, host)
		} else {
			globalConfig = configMap
		}
	}
	// format
	if len(cfg.Hosts) > 0 {
		ioutil.WriteFile(path, []byte(cfg.String()), 0644)
	}
	return hosts, globalConfig
}

func addHost(host *hostConfig) error {
	cfg, aliasMap := parseConfig()
	if err := checkAlias(aliasMap, false, host.aliases); err != nil {
		return err
	}
	if host.config == nil {
		host.config = map[string]string{}
	}
	nodes := []ssh_config.Node{}
	user, hostname, port := parseConnct(host.connect)
	host.connect = formatConnect(user, hostname, port)

	userKV := &ssh_config.KV{Key: USER, Value: user}
	nodes = append(nodes, userKV.SetLeadingSpace(4))

	hostnameKV := &ssh_config.KV{Key: HOSTNAME, Value: hostname}
	nodes = append(nodes, hostnameKV.SetLeadingSpace(4))

	portKV := &ssh_config.KV{Key: PORT, Value: port}
	nodes = append(nodes, portKV.SetLeadingSpace(4))

	checkKeyRepeat := map[string]bool{USER: true, HOSTNAME: true, PORT: true}
	deleteKeys := []string{}
	for k, v := range host.config {
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
		delete(host.config, deleteKey)
	}

	pattern, err := ssh_config.NewPattern(host.aliases)
	if err != nil {
		return nil
	}
	newHost := &ssh_config.Host{
		Patterns: []*ssh_config.Pattern{pattern},
		Nodes:    nodes,
	}
	cfg.Hosts = append(cfg.Hosts, newHost)
	if err := ioutil.WriteFile(path, []byte(cfg.String()), 0644); err != nil {
		return err
	}
	return nil
}

func updateHost(h *hostConfig, newAlias string) error {
	cfg, aliasMap := parseConfig()
	if err := checkAlias(aliasMap, true, h.aliases); err != nil {
		return err
	}

	updateHost := aliasMap[h.aliases]
	if newAlias != "" {
		for _, pattern := range updateHost.Patterns {
			if pattern.String() == h.aliases {
				pattern.SetStr(newAlias)
				h.aliases = newAlias
			}
		}
	}

	updateKV := map[string]string{}
	if h.connect != "" {
		user, hostname, port := parseConnct(h.connect)
		updateKV[USER] = user
		updateKV[HOSTNAME] = hostname
		updateKV[PORT] = port
		h.connect = formatConnect(user, hostname, port)
	}
	if h.config != nil {
		for k, v := range h.config {
			updateKV[strings.ToLower(k)] = v
		}
	}
	h.config = map[string]string{}
	connectMap := map[string]string{USER: "", HOSTNAME: "", PORT: ""}
	for _, node := range updateHost.Nodes {
		switch t := node.(type) {
		case *ssh_config.KV:
			if value, ok := updateKV[t.Key]; ok {
				t.SetLeadingSpace(4)
				t.Value = value
				delete(updateKV, t.Key)
			}
			if _, ok := connectMap[t.Key]; ok {
				connectMap[t.Key] = t.Value
			} else {
				h.config[t.Key] = t.Value
			}
		case *ssh_config.Include:
			// TODO handle include node
		}
	}
	// append new node
	for k, v := range updateKV {
		kv := &ssh_config.KV{Key: k, Value: v}
		updateHost.Nodes = append(updateHost.Nodes, kv.SetLeadingSpace(4))
		if _, ok := connectMap[k]; ok {
			connectMap[k] = v
		} else {
			h.config[k] = v
		}
	}
	h.connect = formatConnect(connectMap[USER], connectMap[HOSTNAME], connectMap[PORT])
	if err := ioutil.WriteFile(path, []byte(cfg.String()), 0644); err != nil {
		return err
	}
	return nil
}

func deleteHost(aliases ...string) error {
	cfg, aliasMap := parseConfig()
	if err := checkAlias(aliasMap, true, aliases...); err != nil {
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
	if err := ioutil.WriteFile(path, []byte(cfg.String()), 0644); err != nil {
		return err
	}
	return nil
}
