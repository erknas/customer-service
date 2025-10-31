package app

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/erknas/customer-service/internal/config"
	"github.com/erknas/customer-service/internal/controller"
	"github.com/erknas/customer-service/internal/repository/postgres"
	"github.com/erknas/customer-service/internal/usecase"
	"github.com/erknas/customer-service/pkg/api/customer"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func Run(cfg *config.Config, logger *zap.Logger) {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	pool, err := pgxpool.New(ctx, cfg.Postgres.URL)
	if err != nil {
		logger.Error("failed to create pgx pool", zap.Error(err))
		os.Exit(1)
	}
	defer pool.Close()

	repo := postgres.New(pool)
	useCase := usecase.New(repo, logger)
	ctrl := controller.New(useCase, logger)

	go runGRPC(cfg, logger, ctrl)
	go runREST(ctx, cfg, logger)

	<-ctx.Done()
}

func runGRPC(cfg *config.Config, logger *zap.Logger, customerServer customer.CustomerServiceServer) {
	port := fmt.Sprintf(":%s", cfg.GRPC.Port)
	ln, err := net.Listen("tcp", port)
	if err != nil {
		logger.Error("falied to listen", zap.Error(err), zap.String("port", port))
		os.Exit(1)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	customer.RegisterCustomerServiceServer(s, customerServer)

	logger.Info("starting gprc server on port", zap.String("port", port))

	if err := s.Serve(ln); err != nil {
		logger.Error("grpc serve failed", zap.Error(err))
	}
}

func runREST(ctx context.Context, cfg *config.Config, logger *zap.Logger) {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	addr := fmt.Sprintf("localhost:%s", cfg.GRPC.Port)
	if err := customer.RegisterCustomerServiceHandlerFromEndpoint(ctx, mux, addr, opts); err != nil {
		logger.Error("failed to register grpc gateway", zap.Error(err))
		os.Exit(1)
	}

	gwPort := fmt.Sprintf(":%s", cfg.GRPC.GatewayPort)

	logger.Info("starting gateway on port", zap.String("port", gwPort))

	if err := http.ListenAndServe(gwPort, mux); err != nil {
		logger.Error("gateway serve failed", zap.Error(err))
	}
}
