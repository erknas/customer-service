package main

import (
	"github.com/erknas/customer-service/internal/app"
	"github.com/erknas/customer-service/internal/config"
	"go.uber.org/zap"
)

func main() {
	cfg := config.Load()

	logger, _ := zap.NewDevelopment()

	app.Run(cfg, logger)
}
