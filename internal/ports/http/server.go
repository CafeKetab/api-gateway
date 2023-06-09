package http

import (
	"encoding/json"
	"fmt"

	"github.com/CafeKetab/gateway/internal/ports/grpc"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type Server struct {
	config *Config
	logger *zap.Logger
	auth   grpc.AuthClient
	app    *fiber.App
}

func New(cfg *Config, log *zap.Logger, auth grpc.AuthClient) *Server {
	server := &Server{config: cfg, logger: log, auth: auth}

	server.app = fiber.New(fiber.Config{JSONEncoder: json.Marshal, JSONDecoder: json.Unmarshal})

	v1 := server.app.Group("/v1")
	v1.Group("/users", server.optionalAuthentication, server.proxy)
	v1.Group("/books", server.requiredAuthentication, server.proxy)

	return server
}

func (server *Server) Serve() error {
	addr := fmt.Sprintf(":%d", server.config.ListenPort)
	if err := server.app.Listen(addr); err != nil {
		server.logger.Error("error resolving server", zap.Error(err))
		return err
	}
	return nil
}
