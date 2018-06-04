package utils

import (
	"errors"
	"os"
	"os/user"
	"sort"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

var (
	SuccessFlag = color.GreenString("✔ ")
	ErrorFlag   = color.RedString("✗ ")
)

// ArgumentsCheck check arguments count correctness
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

// Query values contains keys
func Query(values, keys []string, ignoreCase bool) bool {
	for _, key := range keys {
		if !contains(values, key, ignoreCase) {
			return false
		}
	}
	return true
}

func contains(values []string, key string, ignoreCase bool) bool {
	if ignoreCase {
		key = strings.ToLower(key)
	}
	for _, value := range values {
		if ignoreCase {
			value = strings.ToLower(value)
		}
		if strings.Contains(value, key) {
			return true
		}
	}
	return false
}

// GetHomeDir return user's home directory
func GetHomeDir() string {
	user, err := user.Current()
	if nil == err && user.HomeDir != "" {
		return user.HomeDir
	}
	return os.Getenv("HOME")
}

func GetUsername() string {
	username := ""
	user, err := user.Current()
	if err == nil {
		username = user.Username
	}
	return username
}

func SortKeys(m map[string]string) []string {
	keys := []string{}
	for k, _ := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// ParseConnct parse connect string, format is [user@]host[:port]
func ParseConnct(connect string) (string, string, string) {
	var u, hostname, port string
	hs := strings.SplitN(connect, "@", 2)
	hostname = hs[0]
	if len(hs) == 2 {
		u = hs[0]
		hostname = hs[1]
	}
	hss := strings.SplitN(hostname, ":", 2)
	hostname = hss[0]
	if len(hss) == 2 {
		if _, err := strconv.Atoi(hss[1]); err == nil {
			port = hss[1]
		}
	}
	return u, hostname, port
}
