package server

import (
	"fmt"
	"github.com/bagashiz/go_hexagonal/internal/app/adapters/handlers"
	"github.com/bagashiz/go_hexagonal/internal/app/adapters/storages/db/postgres"
	"github.com/bagashiz/go_hexagonal/internal/app/core/ports"
	"github.com/bagashiz/go_hexagonal/internal/app/infrastructure/configs"
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
