package manssh

import (
	"fmt"

	"github.com/xwjdsh/ssh_config"
)

// HostConfig struct include alias, connect string and other config
type HostConfig struct {
	// Alias alias
	Alias string `json:"alias"`
	// Connection connection
	Connection string `json:"connection"`
	// Path found in which file
	Path string `json:"path"`
	// PathMap key is file path, value is the alias's hosts
	PathMap map[string][]*ssh_config.Host `json:"-"`
	// OwnConfig own config
	OwnConfig map[string]string `json:"own_config,omitempty"`
	// ImplicitConfig implicit config
	ImplicitConfig map[string]string `json:"implicit_config,omitempty"`
}

// NewHostConfig new HostConfig
func NewHostConfig(alias, path string, host *ssh_config.Host) *HostConfig {
	return &HostConfig{
		Alias:          alias,
		Path:           path,
		PathMap:        map[string][]*ssh_config.Host{path: {host}},
		OwnConfig:      map[string]string{},
		ImplicitConfig: map[string]string{},
	}
}

func (hc *HostConfig) connectionStr() string {
	if !hc.Display() {
		return ""
	}

	var (
		user, hostname, port string
		ok                   bool
	)

	if user, ok = hc.OwnConfig["user"]; !ok {
		user = hc.ImplicitConfig["user"]
	}
	if hostname, ok = hc.OwnConfig["hostname"]; !ok {
		hostname = hc.ImplicitConfig["hostname"]
	}
	if port, ok = hc.OwnConfig["port"]; !ok {
		port = hc.ImplicitConfig["port"]
	}

	return fmt.Sprintf("%s@%s:%s", user, hostname, port)
}

// Display Whether to display connection string
func (hc *HostConfig) Display() bool {
	hostname := hc.OwnConfig["hostname"]
	if hostname == "" {
		hostname = hc.ImplicitConfig["hostname"]
	}

	return hostname != ""
}
