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
	// SuccessFlag success flag
	SuccessFlag = color.GreenString("✔ ")
	// ErrorFlag error flag
	ErrorFlag = color.RedString("✗ ")
)

// ArgumentsCheck check arguments count correctness
func ArgumentsCheck(argCount, min, max int) error {
	var err error
	if min > 0 && argCount < min {
		err = errors.New("too few arguments")
	}
	if max > 0 && argCount > max {
		err = errors.New("too many arguments")
	}
	return err
}

// Query values contains keys
func Query(values, keys []string, ignoreCase bool) bool {
	contains := func(key string) bool {
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
	for _, key := range keys {
		if contains(key) {
			return true
		}
	}
	return false
}

// GetHomeDir return user's home directory
func GetHomeDir() string {
	u, err := user.Current()
	if nil == err && u.HomeDir != "" {
		return u.HomeDir
	}
	return os.Getenv("HOME")
}

// GetUsername return current username
func GetUsername() string {
	username := ""
	u, err := user.Current()
	if err == nil {
		username = u.Username
	}
	return username
}

// SortKeys sort map keys
func SortKeys(m map[string]string) []string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// ParseConnect parse connect string, format is [user@]host[:port]
func ParseConnect(connect string) (string, string, string) {
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
