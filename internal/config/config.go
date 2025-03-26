package config

import (
	"mock-api/internal/common/utils"
)

type Config struct {
	Port    string
	MockDir string
}

func NewConfig() *Config {
	return &Config{
		Port:    utils.Env("PORT", "8080"),
		MockDir: utils.Env("MOCK_DIR", "./config/data"),
	}
}
