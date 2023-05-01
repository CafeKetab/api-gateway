package config

import (
	"github.com/CafeKetab/gateway/pkg/logger"
)

type Config struct {
	Logger *logger.Config `koanf:"logger"`
}
