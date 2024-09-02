package config

import "time"

type Config interface {
	String(key string, defaultValue ...string) string
	Int(key string, defaultValue ...int) int
	Duration(key string, defaultValue ...string) time.Duration
	Bool(key string, defaultValue ...bool) bool
}
