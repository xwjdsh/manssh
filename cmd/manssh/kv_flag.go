package main

import (
	"fmt"
	"strings"
)

// https://github.com/urfave/cli/issues/588
type kvFlag struct {
	m map[string]string
}

func (kv *kvFlag) Set(value string) error {
	if value == "" {
		return nil
	}
	if kv.m == nil {
		kv.m = map[string]string{}
	}
	parts := strings.Split(value, "=")
	if len(parts) != 2 || parts[0] == "" {
		return fmt.Errorf("flag param(%s) parse error", value)
	}
	kv.m[parts[0]] = parts[1]
	return nil
}

func (kv *kvFlag) String() string {
	if kv == nil {
		return ""
	}

	item := make([]string, len(kv.m))
	for k, v := range kv.m {
		item = append(item, k+"="+v)
	}
	return strings.Join(item, ",")
}
