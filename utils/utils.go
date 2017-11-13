package utils

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"strconv"
	"strings"

	"github.com/kevinburke/ssh_config"
)

type HostConfig struct {
	Aliases string
	Connect string
	Config  map[string]string
}

func FormatConnect(user, hostname, port string) string {
	return fmt.Sprintf("%s@%s:%s", user, hostname, port)
}

// format is [user@]host[:port]
func ParseConnct(connect string) (string, string, string) {
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

func ArgumentsCheck(argCount, min, max int) error {
	var err error
	if min > 0 && argCount < min {
		err = errors.New("too few arguments.")
	}
	if max > 0 && argCount > max {
		err = errors.New("too many arguments.")
	}
	return err
}

func Query(values, keys []string) bool {
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

func GetHomeDir() string {
	user, err := user.Current()
	if nil == err && user.HomeDir != "" {
		return user.HomeDir
	}
	return os.Getenv("HOME")
}

func CheckAlias(aliasMap map[string]*ssh_config.Host, expectExist bool, aliases ...string) error {
	for _, alias := range aliases {
		ok := aliasMap[alias] != nil
		if !ok && expectExist {
			return fmt.Errorf("alias[%s] not found.", alias)
		} else if ok && !expectExist {
			return fmt.Errorf("alias[%s] already exists.", alias)
		}
	}
	return nil
}
