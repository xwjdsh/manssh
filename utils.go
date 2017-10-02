package main

import (
	"errors"
	"strings"
)

const (
	version = "0.0.1"
	useage  = "management ssh config easier"
)

func argumentsCheck(arguments []string, min, max int) error {
	if len(arguments) < min {
		return errors.New("too few arguments")
	}
	if len(arguments) > max {
		return errors.New("too many arguments")
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
