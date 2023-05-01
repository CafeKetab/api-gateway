package config

import (
	"github.com/CafeKetab/gateway/internal/ports/grpc"
	"github.com/CafeKetab/gateway/internal/ports/http"
	"github.com/CafeKetab/gateway/pkg/logger"
)

func Default() *Config {
	return &Config{
		Logger: &logger.Config{
			Development: true,
			Level:       "debug",
		},
		HTTP: &http.Config{
			ListenPort: 8080,
			TargetUrls: struct {
				Users string "koanf:\"users\""
				Books string "koanf:\"books\""
			}{Users: "localhost:8081", Books: "localhost:8082"},
		},
		GRPC: &grpc.Config{
			AuthGrpcClientAddress: "localhost:9090",
		},
	}
}
