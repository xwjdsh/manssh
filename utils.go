package main

import "strings"

const (
	name    = "manssh"
	version = "0.0.1"
	useage  = "management ssh config easier"
)

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
