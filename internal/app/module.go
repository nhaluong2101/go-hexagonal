package app

import (
	"github.com/bagashiz/go_hexagonal/internal/app/adapters"
	"github.com/bagashiz/go_hexagonal/internal/app/core"
	"github.com/bagashiz/go_hexagonal/internal/app/infrastructure"
	"github.com/bagashiz/go_hexagonal/internal/app/infrastructure/server"
	"go.uber.org/fx"
)

func NewApp() *fx.App {
	return fx.New(
		adapters.Module,
		core.Module,
		infrastructure.Module,
		fx.Invoke(
			server.Serve,
		),
	)
}
