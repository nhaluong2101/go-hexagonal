package server

import (
	"fmt"
	"github.com/bagashiz/go_hexagonal/internal/app/adapters/handlers"
	"github.com/bagashiz/go_hexagonal/internal/app/adapters/storages/db/postgres"
	"github.com/bagashiz/go_hexagonal/internal/app/core/ports"
	"github.com/bagashiz/go_hexagonal/internal/app/infrastructure/configs"
	"go.elastic.co/ecszap"
	"go.uber.org/zap"
	"log/slog"
	"os"
)

// Serve starts the HTTP server
func Serve(
	router *handlers.RouterHandler,
	config *configs.Container,
	db *postgres.DB,
	cache ports.CacheRepository,
) error {

	encoderConfig := ecszap.NewDefaultEncoderConfig()
	core := ecszap.NewCore(encoderConfig, os.Stdout, zap.DebugLevel)
	logger := zap.New(core, zap.AddCaller()).With(zap.String("app", "go-elk")).With(zap.String("environment", "local"))

	logger.Info("application log",
		zap.Int("times", 1),
	)

	// Migrate database
	var err = db.Migrate()
	if err != nil {
		slog.Error("Error migrating database", "error", err)
		os.Exit(1)
	}

	defer cache.Close()
	defer db.Close()

	listenAddr := fmt.Sprintf("%s:%s", config.HTTP.URL, config.HTTP.Port)
	slog.Info("Starting the HTTP server", "listen_address", listenAddr)
	return router.Run(listenAddr)

}
