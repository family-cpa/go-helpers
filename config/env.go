package config

import (
	_ "github.com/joho/godotenv/autoload"
	"os"
	"strconv"
	"time"
)

type EnvConfig struct{}

func NewEnvConfigProvider() Config {
	return &EnvConfig{}
}

func (p *EnvConfig) String(key string, defaultValue ...string) string {
	return defaultString(os.Getenv(key), defaultValue)
}

func (p *EnvConfig) Int(key string, defaultValue ...int) int {
	value, _ := strconv.Atoi(p.String(key))

	if value == 0 && len(defaultValue) > 0 {
		return defaultValue[0]
	}

	return value
}

func (p *EnvConfig) Duration(key string, defaultValue ...string) time.Duration {
	t, _ := time.ParseDuration(defaultString(os.Getenv(key), defaultValue))
	return t
}

func (p *EnvConfig) Bool(key string, defaultValue ...bool) bool {
	value, err := strconv.ParseBool(p.String(key))

	if err != nil && len(defaultValue) > 0 {
		return defaultValue[0]
	}

	return value
}

func defaultString(value string, defaultValue []string) string {
	if len(value) == 0 && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return value
}
