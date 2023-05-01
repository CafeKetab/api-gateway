package config

import (
	"github.com/CafeKetab/gateway/internal/ports/grpc"
	"github.com/CafeKetab/gateway/internal/ports/http"
	"github.com/CafeKetab/gateway/pkg/logger"
)

type Config struct {
	Logger *logger.Config `koanf:"logger"`
	HTTP   *http.Config   `koanf:"http"`
	GRPC   *grpc.Config   `koanf:"grpc"`
}
