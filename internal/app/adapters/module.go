package adapters

import (
	"github.com/bagashiz/go_hexagonal/internal/app/adapters/auth"
	"github.com/bagashiz/go_hexagonal/internal/app/adapters/handlers"
	"github.com/bagashiz/go_hexagonal/internal/app/adapters/storages"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"adapters-module",
	auth.Module,
	storages.Module,
	handlers.Module,
)
