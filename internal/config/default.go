package config

import (
	"github.com/CafeKetab/gateway/pkg/logger"
)

func Default() *Config {
	return &Config{
		Logger: &logger.Config{
			Development: true,
			Level:       "debug",
		},
	}
}
