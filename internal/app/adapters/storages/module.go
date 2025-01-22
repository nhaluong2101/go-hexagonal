package storages

import (
	"github.com/bagashiz/go_hexagonal/internal/app/adapters/repositories"
	"github.com/bagashiz/go_hexagonal/internal/app/adapters/storages/db/postgres"
	"github.com/bagashiz/go_hexagonal/internal/app/adapters/storages/redis"
	"go.uber.org/fx"
)

var Module = fx.Options(
	postgres.Module,
	repositories.Module,
	redis.Module,
)
