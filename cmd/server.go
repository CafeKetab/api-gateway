package cmd

import (
	"os"

	"github.com/CafeKetab/gateway/internal/config"
	"github.com/CafeKetab/gateway/internal/ports/grpc"
	"github.com/CafeKetab/gateway/internal/ports/http"
	"github.com/CafeKetab/gateway/pkg/logger"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type Server struct{}

func (cmd Server) Command(trap chan os.Signal) *cobra.Command {
	run := func(_ *cobra.Command, _ []string) {
		cmd.main(config.Load(true), trap)
	}

	return &cobra.Command{
		Use:   "server",
		Short: "run api-gateway server",
		Run:   run,
	}
}

func (cmd *Server) main(cfg *config.Config, trap chan os.Signal) {
	logger := logger.NewZap(cfg.Logger)

	authGrpcClient := grpc.NewAuthClient(cfg.GRPC, logger)

	httpServer := http.New(cfg.HTTP, logger, authGrpcClient)
	go httpServer.Serve()

	// Keep this at the bottom of the main function
	field := zap.String("signal trap", (<-trap).String())
	logger.Info("exiting by receiving a unix signal", field)
}
